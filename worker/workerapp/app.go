package workerapp

import (
	"github.com/hibiken/asynq"
	"github.com/roessland/withoutings/tasks"
	"github.com/roessland/withoutings/withings"
	"github.com/sirupsen/logrus"
	"log"
	"os"
)

type App struct {
	Log      logrus.FieldLogger
	Withings *withings.Client
	Async    *asynq.Client
}

const redisAddr = "127.0.0.1:6379"

func NewApp() *App {
	app := App{}

	logger := logrus.New()
	app.Log = logger

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

func (app *App) Work() {
	asyncSrv := asynq.NewServer(
		asynq.RedisClientOpt{
			Addr: redisAddr,
		},
		asynq.Config{
			Concurrency: 10,
		},
	)

	mux := asynq.NewServeMux()
	mux.Handle(tasks.TypeWithingsAPICall, tasks.NewWithingsAPICallProcessor())

	if err := asyncSrv.Run(mux); err != nil {
		log.Fatalf("could not run server: %v", err)
	}
}
