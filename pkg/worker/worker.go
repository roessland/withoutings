package worker

import (
	"context"
	"encoding/json"
	"github.com/roessland/withoutings/pkg/logging"
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
	log := logging.MustGetLoggerFromContext(ctx)

	log.WithField("event", "info.worker.started").Info()

	messages, err := wrk.App.Subscriber.Subscribe(ctx, topic.WithingsRawNotificationReceived)
	if err != nil {
		panic(err)
	}

	// TODO refactor to watermill router, avoid acking immediately

	for msg := range messages {
		log = log.WithField("message_uuid", msg.UUID).WithField("message_payload", msg.Payload).WithField("message_metadata", msg.Metadata)
		log.WithField("event", "info.worker.message.received")
		var rawNotificationReceived subscription.RawNotificationReceived
		err := json.Unmarshal(msg.Payload, &rawNotificationReceived)
		if err != nil {
			log.WithError(err).WithField("event", "error.worker.unmarshal.failed").Error()
			continue
		}

		rawNotification, err := wrk.SubscriptionRepo.GetRawNotificationByUUID(ctx,
			rawNotificationReceived.RawNotificationUUID,
		)
		if err != nil {
			log.WithError(err).WithField("event", "error.worker.GetRawNotificationByUUID.failed").Error()
			continue
		}

		err = wrk.Commands.ProcessRawNotification.Handle(ctx, command.ProcessRawNotification{
			RawNotification: rawNotification,
		})
		if err != nil {
			log.WithError(err).WithField("event", "error.worker.ProcessRawNotification.failed").Error()
			continue
		}
		msg.Ack()
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
