package app

import (
	"github.com/alexedwards/scs/v2"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/roessland/withoutings/pkg/config"
	"github.com/roessland/withoutings/pkg/repos/db"
	"github.com/roessland/withoutings/pkg/service/sleep"
	"github.com/roessland/withoutings/pkg/withoutings/app/command"
	"github.com/roessland/withoutings/pkg/withoutings/app/query"
	"github.com/roessland/withoutings/pkg/withoutings/clients/withingsapi"
	"github.com/roessland/withoutings/pkg/withoutings/domain/account"
	"github.com/roessland/withoutings/pkg/withoutings/domain/subscription"
	"github.com/roessland/withoutings/web/templates"
	"github.com/sirupsen/logrus"
)

// App holds all application resources.
type App struct {
	Log              logrus.FieldLogger
	Withings         *withingsapi.Client
	Sessions         *scs.SessionManager
	Templates        templates.Templates
	Sleep            *sleep.Sleep
	DB               *pgxpool.Pool
	Config           *config.Config
	DBQueries        *db.Queries
	AccountRepo      account.Repo
	SubscriptionRepo subscription.Repo
	Commands         Commands
	Queries          Queries
}

type Commands struct {
	SubscribeAccount      command.SubscribeAccountHandler
	CreateOrUpdateAccount command.CreateOrUpdateAccountHandler
}

func (cs Commands) Validate() {
	if cs.SubscribeAccount == nil {
		panic("Commands.SubscribeAccount was nil")
	}
	if cs.CreateOrUpdateAccount == nil {
		panic("Commands.CreateOrUpdateAccount was nil")
	}
}

type Queries struct {
	AccountForWithingsUserID query.AccountByWithingsUserIDHandler
	AccountForUserID         query.AccountByIDHandler
	Accounts                 query.AccountsHandler
}

func (qs Queries) Validate() {
	if qs.AccountForWithingsUserID == nil {
		panic("Queries.AccountForWithingsUserID was nil")
	}

	if qs.AccountForUserID == nil {
		panic("Queries.AccountForUserID was nil")
	}

	if qs.Accounts == nil {
		panic("Queries.Accounts was nil")
	}
}
