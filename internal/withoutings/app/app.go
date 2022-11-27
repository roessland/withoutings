package app

import (
	"github.com/alexedwards/scs/v2"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/roessland/withoutings/internal/config"
	"github.com/roessland/withoutings/internal/repos/db"
	"github.com/roessland/withoutings/internal/service/sleep"
	"github.com/roessland/withoutings/internal/withoutings/adapters/withingsapi"
	"github.com/roessland/withoutings/internal/withoutings/app/command"
	"github.com/roessland/withoutings/internal/withoutings/app/query"
	"github.com/roessland/withoutings/web/templates"
	"github.com/sirupsen/logrus"
)

// App holds all application resources.
type App struct {
	Log              logrus.FieldLogger
	Withings         *withingsapiadapter.Client
	Sessions         *scs.SessionManager
	Templates        templates.Templates
	Sleep            *sleep.Sleep
	DB               *pgxpool.Pool
	Config           *config.Config
	DBQueries        *db.Queries
	AccountRepo      *db.Queries
	SubscriptionRepo *db.Queries
	Commands         Commands
	Queries          Queries
}

type Commands struct {
	CreateOrUpdateAccount command.SubscribeAccountHandler
}

type Queries struct {
	AccountForWithingsUserID query.AccountByWithingsUserIDHandler
}
