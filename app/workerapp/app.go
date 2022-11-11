package workerapp

import (
	"context"
	"fmt"
	"github.com/roessland/withoutings/withings"
	"github.com/sirupsen/logrus"
	"os"
	"time"
)

type WorkerApp struct {
	Log      logrus.FieldLogger
	Withings *withings.Client
	//Async    *asynq.Client
}

const redisAddr = "127.0.0.1:6379"

func NewApp() *WorkerApp {
	app := WorkerApp{}

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

	//app.Async = asynq.NewClient(asynq.RedisClientOpt{
	//	Addr: redisAddr,
	//})

	return &app
}

func (app *WorkerApp) close() {
	//err := app.Async.Close()
	//if err != nil {
	//	app.Log.Print(err)
	//}
}

func (app *WorkerApp) Work(ctx context.Context) {
	//asyncSrv := asynq.NewServer(
	//	asynq.RedisClientOpt{
	//		Addr: redisAddr,
	//	},
	//	asynq.Config{
	//		Concurrency: 10,
	//	},
	//)
	//
	//mux := asynq.NewServeMux()
	//mux.Handle(tasks.TypeWithingsAPICall, tasks.NewWithingsAPICallProcessor())
	//
	//go func() {
	//	if err := asyncSrv.Run(mux); err != nil {
	//		app.Log.Error("could not run server: %v", err)
	//	}
	//}()
	for {

		select {
		case <-ctx.Done():
			app.Log.Info("Worker shutdown initiated.")
			app.close()
			return

		case <-time.After(10 * time.Minute):
			fmt.Println("working")
		}

	}
	//
	//asyncSrv.Shutdown()
	//fmt.Println("WORK SHUTDOWN DONE")

}
