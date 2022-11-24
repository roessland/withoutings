package worker

import (
	"context"
	"fmt"
	"github.com/roessland/withoutings/internal/service"
	"time"
)

type Worker struct {
	*service.App
}

// const redisAddr = "127.0.0.1:6379"

func NewWorker(app *service.App) *Worker {
	return &Worker{app}
}

func (app *Worker) close() {
	//err := app.Async.Close()
	//if err != nil {
	//	app.Log.Print(err)
	//}
}

func (app *Worker) Work(ctx context.Context) {
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
