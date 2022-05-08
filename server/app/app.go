package app

import (
	"github.com/roessland/withoutings/server/domain/services"
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
	Sleep     *services.Sleep
}

func NewApp() *App {
	app := App{}

	logger := logrus.New()
	app.Log = logger

	sessionKey := []byte(os.Getenv("SESSION_KEY"))
	if len(sessionKey) == 0 {
		app.Log.Fatal("SESSION_KEY missing")
	}
	app.Sessions = sessions.NewManager(sessionKey)

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

	app.Sleep = services.NewSleep(app.Withings)

	return &app
}
