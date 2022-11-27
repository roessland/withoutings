package service

import (
	"context"
	"fmt"
	"github.com/alexedwards/scs/pgxstore"
	"github.com/alexedwards/scs/v2"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/roessland/withoutings/internal/config"
	"github.com/roessland/withoutings/internal/repos/db"
	"github.com/roessland/withoutings/internal/service/sleep"
	"github.com/roessland/withoutings/internal/withoutings/adapters"
	"github.com/roessland/withoutings/internal/withoutings/adapters/withingsapi"
	"github.com/roessland/withoutings/internal/withoutings/app"
	"github.com/roessland/withoutings/internal/withoutings/app/command"
	"github.com/roessland/withoutings/internal/withoutings/app/query"
	"github.com/roessland/withoutings/internal/withoutings/domain/account"
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

	svc.DBQueries = db.New(svc.DB)
	svc.AccountRepo = svc.DBQueries

	svc.Sessions = scs.New()
	svc.Sessions.Store = pgxstore.New(svc.DB)

	svc.Withings = withingsapiadapter.NewClient(cfg.WithingsClientID, cfg.WithingsClientSecret, cfg.WithingsRedirectURL)

	svc.Templates = templates.LoadTemplates()

	svc.Sleep = sleep.NewSleep(svc.Withings)

	//
	//svc.Async = asynq.NewClient(asynq.RedisClientOpt{
	//	Addr: redisAddr,
	//})

	return svc, nil
}

func NewApplication(ctx context.Context) app.App {
	return newApplication(ctx)
}

func newApplication(ctx context.Context) app.App {
	cfg, err := config.LoadFromEnv()
	if err != nil {
		panic(fmt.Sprintf("load config: %s", err))
	}

	pool, err := pgxpool.Connect(ctx, cfg.DatabaseURL)
	if err != nil {
		panic(fmt.Sprintf("create connection pool: %s", err))
	}

	queries := db.New(pool)

	var accountRepo account.Repo = adapter.NewAccountPgRepo(queries)

	return app.App{
		Commands: app.Commands{
			CreateOrUpdateAccount: command.NewSubscribeAccountHandler(accountRepo),
		},
		Queries: app.Queries{
			AccountForWithingsUserID: query.NewAccountByWithingsUserIDHandler(accountRepo),
		},
	}
}
