package app

// Inspiration: https://github.com/ThreeDotsLabs/wild-workouts-go-ddd-example/blob/master/internal/trainings/app/app.go

import (
	"context"
	"fmt"
	"github.com/alexedwards/scs/pgxstore"
	"github.com/alexedwards/scs/v2"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/roessland/withoutings/pkg/config"
	"github.com/roessland/withoutings/pkg/db"
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
	"time"
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
}

type MockApp struct {
	*App
	MockWithingsRepo *withings.MockRepo
	MockWithingsSvc  *withingsService.MockService
}

func NewApplication(ctx context.Context, cfg *config.Config) *App {
	logger := logrus.New()

	pool, err := pgxpool.New(ctx, cfg.DatabaseURL)
	if err != nil {
		panic(fmt.Sprintf("create connection pool: %s", err))
	}

	dbQueries := db.New(pool)

	sessions := scs.New()
	sessions.Lifetime = time.Hour * 24 * 180    // 6 months
	sessions.IdleTimeout = time.Hour * 24 * 180 // 6 months

	flashManager := flash.NewManager(sessions)

	sessions.Store = pgxstore.New(pool)

	accountRepo := accountAdapter.NewPgRepo(pool, dbQueries)
	subscriptionRepo := subscriptionAdapter.NewPgRepo(pool, dbQueries)

	withingsHttpClient := withingsAdapter.NewClient(cfg.WithingsClientID, cfg.WithingsClientSecret, cfg.WithingsRedirectURL)

	withingsSvc := withingsService.NewService(withingsHttpClient, accountRepo)

	return &App{
		Log:              logger,
		WithingsRepo:     withingsHttpClient,
		Sessions:         sessions,
		Flash:            flashManager,
		Templates:        templates.NewTemplates(),
		Sleep:            sleep.NewSleep(withingsHttpClient),
		DB:               pool,
		Config:           cfg,
		AccountRepo:      accountRepo,
		SubscriptionRepo: subscriptionRepo,
		WithingsSvc:      withingsSvc,
		Commands: Commands{
			SubscribeAccount:         command.NewSubscribeAccountHandler(accountRepo, subscriptionRepo, withingsSvc, cfg),
			CreateOrUpdateAccount:    command.NewCreateOrUpdateAccountHandler(accountRepo),
			RefreshAccessToken:       command.NewRefreshAccessTokenHandler(accountRepo, withingsHttpClient),
			SyncRevokedSubscriptions: command.NewSyncRevokedSubscriptionsHandler(subscriptionRepo, withingsSvc),
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
