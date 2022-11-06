package webapp

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/roessland/withoutings/domain/services/sleep"
	"github.com/roessland/withoutings/web/sessions"
	"github.com/roessland/withoutings/web/templates"
	"github.com/roessland/withoutings/withings"
	"github.com/sirupsen/logrus"
	"os"
	"time"
)

type WebApp struct {
	Log       logrus.FieldLogger
	Withings  *withings.Client
	Sessions  *sessions.Manager
	Templates templates.Templates
	Sleep     *sleep.Sleep
	//Async     *asynq.Client
	DB *pgxpool.Pool
}

// const redisAddr = "127.0.0.1:6379"

func NewApp(ctx context.Context) *WebApp {
	var err error
	initCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	app := WebApp{}

	logger := logrus.New()
	app.Log = logger

	sessionSecret := []byte(os.Getenv("SESSION_SECRET"))
	if len(sessionSecret) == 0 {
		app.Log.Fatal("SESSION_SECRET missing")
	}
	app.Sessions = sessions.NewManager(sessionSecret)

	withingsClientID := os.Getenv("WITHINGS_CLIENT_ID")
	if withingsClientID == "" {
		app.Log.Fatal("WITHINGS_CLIENT_ID missing")
	}

	withingsClientSecret := os.Getenv("WITHINGS_CLIENT_SECRET")
	if withingsClientSecret == "" {
		app.Log.Fatal("WITHINGS_CLIENT_SECRET missing")
	}

	withingsRedirectURL := os.Getenv("WITHINGS_REDIRECT_URL")
	if withingsRedirectURL == "" {
		app.Log.Fatal("WITHINGS_REDIRECT_URL missing")
	}

	dbConnectionString := os.Getenv("WOT_DATABASE_URL")
	if dbConnectionString == "" {
		app.Log.Fatal("WOT_DATABASE_URL missing")
	}
	app.DB, err = pgxpool.New(initCtx, dbConnectionString)
	if err != nil {
		app.Log.Fatalf("Unable to create connection pool: %v", err)
	}

	app.Withings = withings.NewClient(withingsClientID, withingsClientSecret, withingsRedirectURL)

	app.Templates = templates.LoadTemplates()

	app.Sleep = sleep.NewSleep(app.Withings)
	//
	//app.Async = asynq.NewClient(asynq.RedisClientOpt{
	//	Addr: redisAddr,
	//})

	return &app
}

func (app *WebApp) Close() {
	//err := app.Async.Close()
	//if err != nil {
	//	app.Log.Print(err)
	//}

	app.DB.Close()
}
