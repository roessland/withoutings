package app

// Inspiration: https://github.com/ThreeDotsLabs/wild-workouts-go-ddd-example/blob/master/internal/trainings/app/app.go

import (
	"github.com/alexedwards/scs/v2"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/roessland/withoutings/pkg/config"
	"github.com/roessland/withoutings/pkg/db"
	"github.com/roessland/withoutings/pkg/service/sleep"
	"github.com/roessland/withoutings/pkg/withoutings/app/command"
	"github.com/roessland/withoutings/pkg/withoutings/app/query"
	"github.com/roessland/withoutings/pkg/withoutings/domain/account"
	"github.com/roessland/withoutings/pkg/withoutings/domain/subscription"
	"github.com/roessland/withoutings/pkg/withoutings/domain/withings"
	"github.com/roessland/withoutings/web/templates"
	"github.com/sirupsen/logrus"
)

// App holds all application resources.
type App struct {
	Log              logrus.FieldLogger
	Sessions         *scs.SessionManager
	Templates        templates.Templates
	Sleep            *sleep.Sleep
	DB               *pgxpool.Pool
	Config           *config.Config
	DBQueries        *db.Queries
	WithingsRepo     withings.Repo
	AccountRepo      account.Repo
	SubscriptionRepo subscription.Repo
	Commands         Commands
	Queries          Queries
}

type Commands struct {
	SubscribeAccount      command.SubscribeAccountHandler
	CreateOrUpdateAccount command.CreateOrUpdateAccountHandler
	RefreshAccessToken    command.RefreshAccessTokenHandler
}

func (cs Commands) Validate() {
	if cs.SubscribeAccount == nil {
		panic("Commands.SubscribeAccount was nil")
	}
	if cs.CreateOrUpdateAccount == nil {
		panic("Commands.CreateOrUpdateAccount was nil")
	}
	if cs.RefreshAccessToken == nil {
		panic("Commands.RefreshAccessToken was nil")
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
