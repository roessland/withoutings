package service

import (
	"context"
	"fmt"
	"github.com/alexedwards/scs/pgxstore"
	"github.com/alexedwards/scs/v2"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/roessland/withoutings/pkg/config"
	"github.com/roessland/withoutings/pkg/db"
	"github.com/roessland/withoutings/pkg/service/sleep"
	"github.com/roessland/withoutings/pkg/withoutings/adapter/account"
	"github.com/roessland/withoutings/pkg/withoutings/adapter/subscription"
	withingsAdapter "github.com/roessland/withoutings/pkg/withoutings/adapter/withings"
	"github.com/roessland/withoutings/pkg/withoutings/app"
	"github.com/roessland/withoutings/pkg/withoutings/app/command"
	"github.com/roessland/withoutings/pkg/withoutings/app/query"
	"github.com/roessland/withoutings/web/templates"
	"github.com/sirupsen/logrus"
	"time"
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
	sessions.Lifetime = time.Hour * 24 * 180    // 6 months
	sessions.IdleTimeout = time.Hour * 24 * 180 // 6 months

	sessions.Store = pgxstore.New(pool)

	accountRepo := account.NewPgRepo(pool, dbQueries)
	subscriptionRepo := subscription.NewPgRepo(pool, dbQueries)

	withingsHttpClient := withingsAdapter.NewClient(cfg.WithingsClientID, cfg.WithingsClientSecret, cfg.WithingsRedirectURL)

	return &app.App{
		Log:              logger,
		WithingsRepo:     withingsHttpClient,
		Sessions:         sessions,
		Templates:        templates.LoadTemplates(),
		Sleep:            sleep.NewSleep(withingsHttpClient),
		DB:               pool,
		Config:           cfg,
		DBQueries:        dbQueries,
		AccountRepo:      accountRepo,
		SubscriptionRepo: subscriptionRepo,
		Commands: app.Commands{
			SubscribeAccount:      command.NewSubscribeAccountHandler(accountRepo, subscriptionRepo, withingsHttpClient, cfg),
			CreateOrUpdateAccount: command.NewCreateOrUpdateAccountHandler(accountRepo),
			RefreshAccessToken:    command.NewRefreshAccessTokenHandler(accountRepo, withingsHttpClient),
		},
		Queries: app.Queries{
			AccountByWithingsUserID: query.NewAccountByWithingsUserIDHandler(accountRepo),
			AccountByAccountUUID:    query.NewAccountByUUIDHandler(accountRepo),
			AllAccounts:             query.NewAllAccountsHandler(accountRepo),
		},
	}
}
