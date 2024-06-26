package app

import (
	"context"
	"database/sql"
	wmSql "github.com/ThreeDotsLabs/watermill-sql/v3/pkg/sql"
	"testing"

	"github.com/alexedwards/scs/pgxstore"
	"github.com/alexedwards/scs/v2"
	"github.com/alexedwards/scs/v2/memstore"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/roessland/withoutings/pkg/config"
	"github.com/roessland/withoutings/pkg/db"
	"github.com/roessland/withoutings/pkg/logging"
	"github.com/roessland/withoutings/pkg/service/sleep"
	"github.com/roessland/withoutings/pkg/testctx"
	"github.com/roessland/withoutings/pkg/web/flash"
	"github.com/roessland/withoutings/pkg/web/templates"
	accountAdapter "github.com/roessland/withoutings/pkg/withoutings/adapter/account"
	subscriptionAdapter "github.com/roessland/withoutings/pkg/withoutings/adapter/subscription"
	"github.com/roessland/withoutings/pkg/withoutings/app/command"
	"github.com/roessland/withoutings/pkg/withoutings/app/query"
	withingsService "github.com/roessland/withoutings/pkg/withoutings/app/service/withings"
	"github.com/roessland/withoutings/pkg/withoutings/domain/account"
	"github.com/roessland/withoutings/pkg/withoutings/domain/subscription"
	"github.com/roessland/withoutings/pkg/withoutings/domain/withings"
)

func NewTestApplication(t *testing.T, ctx context.Context, database *pgxpool.Pool, stdDB *sql.DB) *MockApp {
	// TODO pass in config, instead of database connections.
	// TODO: Merge with NewApplication, avoid code duplication
	cfg := &config.Config{
		ListenAddr:            "<test-listen-addr>",
		SessionSecret:         []byte("abc123"),
		WebsiteURL:            "https://withoutings.com/",
		WithingsClientID:      "the_client_id",
		WithingsClientSecret:  "the_client_secret",
		WithingsRedirectURL:   "https://withoutings.com/auth/callback",
		WithingsWebhookSecret: "qwerty1234",
		DatabaseURL:           "<test-database-url>",
	}
	log := logging.MustGetLoggerFromContext(ctx)

	dbQueries := db.New(database)
	accountRepo := accountAdapter.NewPgRepo(database, dbQueries)
	subscriptionRepo := subscriptionAdapter.NewPgRepo(database, dbQueries)
	mockWithingsRepo := withings.NewMockRepo(t)
	mockWithingsSvc := withingsService.NewMockService(t)

	// Watermill SQL PubSub
	watermillLogger := logging.NewLogrusWatermill(log)
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

	queries := Queries{
		AccountByAccountUUID:    query.NewAccountByUUIDHandler(accountRepo),
		AccountByWithingsUserID: query.NewAccountByWithingsUserIDHandler(accountRepo),
		AllAccounts:             query.NewAllAccountsHandler(accountRepo),
	}
	commands := Commands{
		SubscribeAccount:         command.NewSubscribeAccountHandler(accountRepo, subscriptionRepo, mockWithingsSvc, cfg),
		CreateOrUpdateAccount:    command.NewCreateOrUpdateAccountHandler(accountRepo),
		RefreshAccessToken:       command.NewRefreshAccessTokenHandler(accountRepo, mockWithingsRepo),
		SyncRevokedSubscriptions: command.NewSyncRevokedSubscriptionsHandler(subscriptionRepo, mockWithingsSvc),
		ProcessRawNotification:   command.NewProcessRawNotificationHandler(subscriptionRepo, mockWithingsSvc, accountRepo),
		FetchNotificationData:    command.NewFetchNotificationDataHandler(subscriptionRepo, mockWithingsSvc, accountRepo, publisher),
	}

	sleepSvc := sleep.NewSleep(nil) // no http client for now

	templateSvc := templates.NewTemplates(templates.Config{})

	sessionManager := scs.New()
	sessionManager.Store = pgxstore.New(database)

	flashManager := flash.NewManager(sessionManager)

	svc := &MockApp{
		App: &App{
			Log:              log,
			Sessions:         sessionManager,
			Flash:            flashManager,
			Templates:        templateSvc,
			Sleep:            sleepSvc,
			DB:               database,
			Config:           cfg,
			WithingsRepo:     mockWithingsRepo,
			AccountRepo:      accountRepo,
			SubscriptionRepo: subscriptionRepo,
			WithingsSvc:      mockWithingsSvc,
			Commands:         commands,
			Queries:          queries,
			Publisher:        publisher,
			Subscriber:       subscriber,
		},
		MockWithingsRepo: mockWithingsRepo,
		MockWithingsSvc:  mockWithingsSvc,
	}
	return svc
}

// NewMockApplication returns a bare-bones mock application that can't do much.
// Use NewTestApplication for integration tests.
func NewMockApplication(t *testing.T) *App {
	ctx := testctx.New()
	svc := &App{}
	svc.Log = ctx.Logger
	svc.Sessions = newInMemorySessionsManager()
	svc.Flash = flash.NewManager(svc.Sessions)
	svc.Templates = templates.NewTemplates(templates.Config{})
	svc.Sleep = sleep.NewSleep(nil)
	svc.Config = &config.Config{}
	svc.WithingsRepo = withings.NewMockRepo(t)
	svc.AccountRepo = account.NewMockRepo(t)
	svc.SubscriptionRepo = subscription.NewMockRepo(t)
	svc.WithingsSvc = withingsService.NewMockService(t)
	svc.Commands = Commands{
		SubscribeAccount:         command.NewSubscribeAccountHandler(svc.AccountRepo, svc.SubscriptionRepo, svc.WithingsSvc, svc.Config),
		CreateOrUpdateAccount:    command.NewCreateOrUpdateAccountHandler(svc.AccountRepo),
		RefreshAccessToken:       command.NewRefreshAccessTokenHandler(svc.AccountRepo, svc.WithingsRepo),
		SyncRevokedSubscriptions: command.NewSyncRevokedSubscriptionsHandler(svc.SubscriptionRepo, svc.WithingsSvc),
		ProcessRawNotification:   command.NewProcessRawNotificationHandler(svc.SubscriptionRepo, svc.WithingsSvc, svc.AccountRepo),
		FetchNotificationData:    command.NewFetchNotificationDataHandler(svc.SubscriptionRepo, svc.WithingsSvc, svc.AccountRepo, svc.Publisher),
	}
	svc.Queries = Queries{
		AccountByWithingsUserID: query.NewAccountByWithingsUserIDHandler(svc.AccountRepo),
		AccountByAccountUUID:    query.NewAccountByUUIDHandler(svc.AccountRepo),
		AllAccounts:             query.NewAllAccountsHandler(svc.AccountRepo),
	}
	svc.Validate()
	return svc
}

func newInMemorySessionsManager() *scs.SessionManager {
	sessionManager := scs.New()
	sessionManager.Store = memstore.New()
	return sessionManager
}
