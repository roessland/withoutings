package service

import (
	"context"
	"fmt"
	"github.com/alexedwards/scs/pgxstore"
	"github.com/alexedwards/scs/v2"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/roessland/withoutings/pkg/config"
	"github.com/roessland/withoutings/pkg/repos/db"
	"github.com/roessland/withoutings/pkg/service/sleep"
	"github.com/roessland/withoutings/pkg/withoutings/adapter"
	"github.com/roessland/withoutings/pkg/withoutings/app"
	"github.com/roessland/withoutings/pkg/withoutings/app/command"
	"github.com/roessland/withoutings/pkg/withoutings/app/query"
	"github.com/roessland/withoutings/pkg/withoutings/clients/withingsapi"
	"github.com/roessland/withoutings/web/templates"
	"github.com/sirupsen/logrus"
	"time"
)

// const redisAddr = "127.0.0.1:6379"

// NewApplication creates a new Withoutings service.
func NewApplication_(ctx context.Context) (*app.App, error) {
	svc := &app.App{}

	var err error

	initCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	logger := logrus.New()
	svc.Log = logger

	cfg, err := config.LoadFromEnv()
	if err != nil {
		return svc, fmt.Errorf("load config: %w", err)
	}
	svc.Config = cfg

	svc.DB, err = pgxpool.Connect(initCtx, cfg.DatabaseURL)
	if err != nil {
		return svc, fmt.Errorf("create connection pool: %w", err)
	}

	queries := db.New(svc.DB)
	svc.DBQueries = queries
	svc.AccountRepo = adapter.NewAccountPgRepo(queries)

	svc.Sessions = scs.New()
	svc.Sessions.Store = pgxstore.New(svc.DB)

	svc.Withings = withingsapi.NewClient(cfg.WithingsClientID, cfg.WithingsClientSecret, cfg.WithingsRedirectURL)

	svc.Templates = templates.LoadTemplates()

	svc.Sleep = sleep.NewSleep(svc.Withings)

	//
	//svc.Async = asynq.NewClient(asynq.RedisClientOpt{
	//	Addr: redisAddr,
	//})

	return svc, nil
}

func NewApplication(ctx context.Context) *app.App {
	return newApplication(ctx)
}

func newApplication(ctx context.Context) *app.App {
	cfg, err := config.LoadFromEnv()
	if err != nil {
		panic(fmt.Sprintf("load config: %s", err))
	}

	logger := logrus.New()

	pool, err := pgxpool.Connect(ctx, cfg.DatabaseURL)
	if err != nil {
		panic(fmt.Sprintf("create connection pool: %s", err))
	}

	dbQueries := db.New(pool)

	sessions := scs.New()
	sessions.Store = pgxstore.New(pool)

	accountRepo := adapter.NewAccountPgRepo(dbQueries)
	subscriptionRepo := dbQueries // TODO refactor to adapter

	withings := withingsapi.NewClient(cfg.WithingsClientID, cfg.WithingsClientSecret, cfg.WithingsRedirectURL)

	return &app.App{
		Log:              logger,
		Withings:         withings,
		Sessions:         sessions,
		Templates:        templates.LoadTemplates(),
		Sleep:            sleep.NewSleep(withings),
		DB:               pool,
		Config:           cfg,
		DBQueries:        dbQueries,
		AccountRepo:      accountRepo,
		SubscriptionRepo: subscriptionRepo,
		Commands: app.Commands{
			SubscribeAccount:      command.NewSubscribeAccountHandler(accountRepo),
			CreateOrUpdateAccount: command.NewCreateOrUpdateAccountHandler(accountRepo),
		},
		Queries: app.Queries{
			AccountForWithingsUserID: query.NewAccountByWithingsUserIDHandler(accountRepo),
			AccountForUserID:         query.NewAccountByIDHandler(accountRepo),
			Accounts:                 query.NewAccountsHandler(accountRepo),
		},
	}
}
