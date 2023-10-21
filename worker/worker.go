package worker

import (
	"context"
	"fmt"
	"time"

	"github.com/roessland/withoutings/pkg/withoutings/app"
)

type Worker struct {
	*app.App
}

// const redisAddr = "127.0.0.1:6379"

func NewWorker(svc *app.App) *Worker {
	return &Worker{svc}
}

func (wrk *Worker) close() {
	//err := app.Async.Close()
	//if err != nil {
	//	app.Log.Print(err)
	//}
}

func (wrk *Worker) Work(ctx context.Context) {
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
			wrk.Log.WithField("event", "worker.shutdown.initiated").Info()
			wrk.close()
			return

		case <-time.After(10 * time.Minute):
			fmt.Println("working")
		}
	}
	//
	//asyncSrv.Shutdown()
	//fmt.Println("WORK SHUTDOWN DONE")
}
