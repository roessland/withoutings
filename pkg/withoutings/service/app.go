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
)

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

	accountRepo := adapter.NewAccountPgRepo(pool, dbQueries)
	subscriptionRepo := adapter.NewSubscriptionPgRepo(pool, dbQueries)

	withingsClient := withingsapi.NewClient(cfg.WithingsClientID, cfg.WithingsClientSecret, cfg.WithingsRedirectURL)

	return &app.App{
		Log:              logger,
		Withings:         withingsClient,
		Sessions:         sessions,
		Templates:        templates.LoadTemplates(),
		Sleep:            sleep.NewSleep(withingsClient),
		DB:               pool,
		Config:           cfg,
		DBQueries:        dbQueries,
		AccountRepo:      accountRepo,
		SubscriptionRepo: subscriptionRepo,
		Commands: app.Commands{
			SubscribeAccount:      command.NewSubscribeAccountHandler(accountRepo, subscriptionRepo, withingsClient),
			CreateOrUpdateAccount: command.NewCreateOrUpdateAccountHandler(accountRepo),
		},
		Queries: app.Queries{
			AccountForWithingsUserID: query.NewAccountByWithingsUserIDHandler(accountRepo),
			AccountForUserID:         query.NewAccountByIDHandler(accountRepo),
			Accounts:                 query.NewAccountsHandler(accountRepo),
		},
	}
}
