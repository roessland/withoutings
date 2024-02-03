package app

// Inspiration: https://github.com/ThreeDotsLabs/wild-workouts-go-ddd-example/blob/master/internal/trainings/app/app.go

import (
	"context"
	"database/sql"
	"fmt"
	wmSql "github.com/ThreeDotsLabs/watermill-sql/v2/pkg/sql"
	"github.com/ThreeDotsLabs/watermill/message"
	"time"

	"github.com/alexedwards/scs/postgresstore"
	"github.com/alexedwards/scs/v2"
	"github.com/alexedwards/scs/v2/memstore"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/roessland/withoutings/pkg/config"
	"github.com/roessland/withoutings/pkg/db"
	"github.com/roessland/withoutings/pkg/logging"
	"github.com/roessland/withoutings/pkg/service/sleep"
	"github.com/roessland/withoutings/pkg/web/flash"
	"github.com/roessland/withoutings/pkg/web/templates"
	accountAdapter "github.com/roessland/withoutings/pkg/withoutings/adapter/account"
	subscriptionAdapter "github.com/roessland/withoutings/pkg/withoutings/adapter/subscription"
	withingsAdapter "github.com/roessland/withoutings/pkg/withoutings/adapter/withings"
	"github.com/roessland/withoutings/pkg/withoutings/app/command"
	"github.com/roessland/withoutings/pkg/withoutings/app/query"
	withingsService "github.com/roessland/withoutings/pkg/withoutings/app/service/withings"
	"github.com/roessland/withoutings/pkg/withoutings/domain/account"
	"github.com/roessland/withoutings/pkg/withoutings/domain/subscription"
	"github.com/roessland/withoutings/pkg/withoutings/domain/withings"

	"github.com/sirupsen/logrus"
)

// App holds all application resources.
type App struct {
	Log              logrus.FieldLogger
	Sessions         *scs.SessionManager
	Flash            *flash.Manager
	Templates        *templates.Templates
	Sleep            *sleep.Sleep
	DB               *pgxpool.Pool
	Config           *config.Config
	WithingsRepo     withings.Repo
	AccountRepo      account.Repo
	WithingsSvc      withingsService.Service
	SubscriptionRepo subscription.Repo
	Commands         Commands
	Queries          Queries
	Publisher        message.Publisher
	Subscriber       message.Subscriber
}

type MockApp struct {
	*App
	MockWithingsRepo *withings.MockRepo
	MockWithingsSvc  *withingsService.MockService
}

func NewApplication(ctx context.Context, cfg *config.Config) *App {
	log := logging.NewLogger(cfg.LogFormat)

	// Two separate pools, because the watermill SQL PubSub
	// requires a *sql.DB, while the rest of the application uses a *pgxpool.Pool.
	stdDB, err := sql.Open("pgx", cfg.DatabaseURL)
	if err != nil {
		panic(fmt.Sprintf("create DB connection: %s", err))
	}

	pool, err := pgxpool.New(ctx, cfg.DatabaseURL)
	if err != nil {
		panic(fmt.Sprintf("create DB connection pool: %s", err))
	}
	go func() {
		<-ctx.Done()
		stdDB.Close()
		pool.Close()
	}()

	dbQueries := db.New(pool)

	sessions := scs.New()
	sessions.Lifetime = time.Hour * 24 * 180    // 6 months
	sessions.IdleTimeout = time.Hour * 24 * 180 // 6 months
	go func(store scs.Store) {
		// Stop useless memstore cleanup goroutine started by scs.New()
		if m, ok := store.(*memstore.MemStore); ok {
			time.Sleep(time.Millisecond)
			m.StopCleanup()
		}
	}(sessions.Store)
	pgSessionsStore := postgresstore.New(stdDB)
	go func() {
		<-ctx.Done()
		pgSessionsStore.StopCleanup()
	}()
	sessions.Store = pgSessionsStore

	flashManager := flash.NewManager(sessions)

	accountRepo := accountAdapter.NewPgRepo(pool, dbQueries)
	subscriptionRepo := subscriptionAdapter.NewPgRepo(pool, dbQueries)

	withingsHttpClient := withingsAdapter.NewClient(cfg.WithingsClientID, cfg.WithingsClientSecret, cfg.WithingsRedirectURL)

	withingsSvc := withingsService.NewService(withingsHttpClient, accountRepo)

	templateSvc := templates.NewTemplates(templates.Config{})
	log.
		WithField("event", "templates.loaded").
		WithField("template-source", templateSvc.Source()).
		Info()

	watermillLogger := logging.NewLogrusWatermill(log)

	// Watermill SQL PubSub
	sqlPublisher, err := wmSql.NewPublisher(
		stdDB,
		wmSql.PublisherConfig{
			SchemaAdapter: wmSql.DefaultPostgreSQLSchema{},
		},
		watermillLogger,
	)
	if err != nil {
		panic(err)
	}
	sqlSubscriber, err := wmSql.NewSubscriber(
		stdDB,
		wmSql.SubscriberConfig{
			SchemaAdapter:  wmSql.DefaultPostgreSQLSchema{},
			OffsetsAdapter: wmSql.DefaultPostgreSQLOffsetsAdapter{},
		},
		watermillLogger,
	)
	if err != nil {
		panic(err)
	}
	publisher := sqlPublisher
	subscriber := sqlSubscriber

	return &App{
		Log:              log,
		WithingsRepo:     withingsHttpClient,
		Sessions:         sessions,
		Flash:            flashManager,
		Templates:        templateSvc,
		Sleep:            sleep.NewSleep(withingsSvc),
		DB:               pool,
		Config:           cfg,
		AccountRepo:      accountRepo,
		SubscriptionRepo: subscriptionRepo,
		WithingsSvc:      withingsSvc,
		Publisher:        publisher,
		Subscriber:       subscriber,
		Commands: Commands{
			SubscribeAccount:         command.NewSubscribeAccountHandler(accountRepo, subscriptionRepo, withingsSvc, cfg),
			CreateOrUpdateAccount:    command.NewCreateOrUpdateAccountHandler(accountRepo),
			RefreshAccessToken:       command.NewRefreshAccessTokenHandler(accountRepo, withingsHttpClient),
			SyncRevokedSubscriptions: command.NewSyncRevokedSubscriptionsHandler(subscriptionRepo, withingsSvc),
			ProcessRawNotification:   command.NewProcessRawNotificationHandler(subscriptionRepo, withingsSvc, accountRepo),
			FetchNotificationData:    command.NewFetchNotificationDataHandler(subscriptionRepo, withingsSvc, accountRepo, publisher),
		},
		Queries: Queries{
			AccountByWithingsUserID: query.NewAccountByWithingsUserIDHandler(accountRepo),
			AccountByAccountUUID:    query.NewAccountByUUIDHandler(accountRepo),
			AllAccounts:             query.NewAllAccountsHandler(accountRepo),
		},
	}
}

func (svc *App) Validate() {
	if svc.Log == nil {
		panic("App.Log was nil")
	}
	if svc.Sessions == nil {
		panic("App.Sessions was nil")
	}
	if svc.Flash == nil {
		panic("App.Flash was nil")
	}
	if svc.Templates == nil {
		panic("App.Templates was nil")
	}
	if svc.Sleep == nil {
		panic("App.Sleep was nil")
	}
	if svc.Config == nil {
		panic("App.Config was nil")
	}
	if svc.WithingsRepo == nil {
		panic("App.WithingsRepo was nil")
	}
	if svc.AccountRepo == nil {
		panic("App.AccountRepo was nil")
	}
	if svc.SubscriptionRepo == nil {
		panic("App.SubscriptionRepo was nil")
	}
	svc.Commands.Validate()
	svc.Queries.Validate()
}

type Commands struct {
	SubscribeAccount         command.SubscribeAccountHandler
	CreateOrUpdateAccount    command.CreateOrUpdateAccountHandler
	RefreshAccessToken       command.RefreshAccessTokenHandler
	SyncRevokedSubscriptions command.SyncRevokedSubscriptionsHandler
	ProcessRawNotification   command.ProcessRawNotificationHandler
	FetchNotificationData    command.FetchNotificationDataHandler
}

func (cs Commands) Validate() {
	if cs.SubscribeAccount == nil {
		panic("Commands.SyncRevokedSubscriptions was nil")
	}
	if cs.CreateOrUpdateAccount == nil {
		panic("Commands.CreateOrUpdateAccount was nil")
	}
	if cs.RefreshAccessToken == nil {
		panic("Commands.RefreshAccessToken was nil")
	}
	if cs.SyncRevokedSubscriptions == nil {
		panic("Commands.SyncRevokedSubscriptions was nil")
	}
	if cs.ProcessRawNotification == nil {
		panic("Commands.ProcessRawNotification was nil")
	}
	if cs.FetchNotificationData == nil {
		panic("Commands.FetchNotificationData was nil")
	}
}

type Queries struct {
	AccountByWithingsUserID query.AccountByWithingsUserIDHandler
	AccountByAccountUUID    query.AccountByUUIDHandler
	AllAccounts             query.AllAccountsHandler
}

func (qs Queries) Validate() {
	if qs.AccountByWithingsUserID == nil {
		panic("Queries.AccountByWithingsUserID was nil")
	}

	if qs.AccountByAccountUUID == nil {
		panic("Queries.AccountByAccountUUID was nil")
	}

	if qs.AllAccounts == nil {
		panic("Queries.AllAccounts was nil")
	}
}
