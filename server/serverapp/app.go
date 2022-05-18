package serverapp

import (
	"github.com/hibiken/asynq"
	"github.com/roessland/withoutings/domain/services/sleep"
	"github.com/roessland/withoutings/server/sessions"
	"github.com/roessland/withoutings/server/templates"
	"github.com/roessland/withoutings/withings"
	"github.com/sirupsen/logrus"
	"os"
)

type App struct {
	Log       logrus.FieldLogger
	Withings  *withings.Client
	Sessions  *sessions.Manager
	Templates templates.Templates
	Sleep     *sleep.Sleep
	Async     *asynq.Client
}

const redisAddr = "127.0.0.1:6379"

func NewApp() *App {
	app := App{}

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

	app.Withings = withings.NewClient(withingsClientID, withingsClientSecret, withingsRedirectURL)

	app.Templates = templates.LoadTemplates()

	app.Sleep = sleep.NewSleep(app.Withings)

	app.Async = asynq.NewClient(asynq.RedisClientOpt{
		Addr: redisAddr,
	})

	return &app
}

func (app *App) Close() {
	err := app.Async.Close()
	if err != nil {
		app.Log.Print(err)
	}
}
