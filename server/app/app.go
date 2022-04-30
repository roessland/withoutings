package app

import (
	"github.com/gorilla/sessions"
	"github.com/roessland/withoutings/withings"
	"github.com/sirupsen/logrus"
	"os"
)

type App struct {
	Log            logrus.FieldLogger
	WithingsClient *withings.Client
	SessionKey     []byte
	CookieStore    *sessions.CookieStore
}

func NewApp() *App {
	app := App{}

	logger := logrus.New()
	app.Log = logger

	app.SessionKey = []byte(os.Getenv("SESSION_KEY"))
	if len(app.SessionKey) == 0 {
		app.Log.Fatal("SESSION_KEY missing")
	}
	app.CookieStore = sessions.NewCookieStore(app.SessionKey)

	withingsClientID := os.Getenv("WITHINGS_CLIENT_ID")
	if withingsClientID == "" {
		app.Log.Fatal("WITHINGS_CLIENT_ID missing")
	}

	withingsClientSecret := os.Getenv("WITHINGS_CLIENT_SECRET")
	if withingsClientSecret == "" {
		app.Log.Fatal("WITHINGS_CLIENT_SECRE missing")
	}

	withingsRedirectURL := os.Getenv("WITHINGS_REDIRECT_URL")
	if withingsRedirectURL == "" {
		app.Log.Fatal("WITHINGS_REDIRECT_URL missing")
	}

	app.WithingsClient = withings.NewClient(withingsClientID, withingsClientSecret, withingsRedirectURL)

	return &app
}
