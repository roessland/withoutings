package worker

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/roessland/withoutings/pkg/withoutings/adapter/topic"
	"github.com/roessland/withoutings/pkg/withoutings/app"
	"github.com/roessland/withoutings/pkg/withoutings/app/command"
	"github.com/roessland/withoutings/pkg/withoutings/domain/subscription"
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

	messages, err := wrk.App.Subscriber.Subscribe(ctx, topic.WithingsRawNotificationReceived)
	if err != nil {
		panic(err)
	}

	// TODO refactor to watermill router, avoid acking immediately

	for msg := range messages {
		fmt.Printf("Received message: %s - %s\n", msg.UUID, msg.Payload)
		var rawNotificationReceived subscription.RawNotificationReceived
		err := json.Unmarshal(msg.Payload, &rawNotificationReceived)
		if err != nil {
			panic(err)
		}

		rawNotification, err := wrk.SubscriptionRepo.GetRawNotificationByUUID(ctx,
			rawNotificationReceived.RawNotificationUUID,
		)
		if err != nil {
			panic(err)
		}

		err = wrk.Commands.ProcessRawNotification.Handle(ctx, command.ProcessRawNotification{
			RawNotification: rawNotification,
		})
		if err != nil {
			panic(err)
		}
	}
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
	//for {
	//	select {
	//	case <-ctx.Done():
	//		wrk.Log.WithField("event", "info.worker.shutdown.initiated").Info()
	//		wrk.close()
	//		return
	//
	//	case <-time.After(10 * time.Minute):
	//		wrk.Log.WithField("event", "info.worker.heartbeat").Info()
	//	}
	//}
	//
	//asyncSrv.Shutdown()
	//fmt.Println("WORK SHUTDOWN DONE")
}
