package worker

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/google/uuid"
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
	wmLog := logging.NewLogrusWatermill(log)

	log.WithField("event", "info.worker.started").Info()

	router, err := message.NewRouter(message.RouterConfig{}, wmLog)
	if err != nil {
		log.WithError(err).WithField("event", "panic.worker.NewRouter.failed").Error()
		msg := fmt.Errorf("unable to create router: %w", err)
		panic(msg)
	}

	// TODO: Refactor into struct handler
	router.AddHandler(
		"process_raw_notification",
		topic.WithingsRawNotificationReceived,
		wrk.Subscriber,
		topic.WithingsNotificationReceived,
		wrk.Publisher,
		func(msg *message.Message) ([]*message.Message, error) {
			log = log.
				WithField("handler_name", "process_raw_notification").
				WithField("message_uuid", msg.UUID).
				WithField("message_payload", string(msg.Payload)).
				WithField("message_metadata", msg.Metadata)

			log.WithField("event", "info.msghandler.process_raw_notification.started").Info()

			var rawNotificationReceived subscription.RawNotificationReceived
			err := json.Unmarshal(msg.Payload, &rawNotificationReceived)
			if err != nil {
				log.WithError(err).WithField("event", "error.msghandler.process_raw_notification.unmarshal.failed").Error()
				return nil, err
			}

			rawNotification, err := wrk.SubscriptionRepo.GetRawNotificationByUUID(ctx,
				rawNotificationReceived.RawNotificationUUID,
			)
			if err != nil {
				log.WithError(err).WithField("event", "error.msghandler.process_raw_notification.GetRawNotificationByUUID.failed").Error()
				return nil, err
			}

			err = wrk.Commands.ProcessRawNotification.Handle(ctx, command.ProcessRawNotification{
				RawNotification: rawNotification,
			})
			if err != nil {
				log.WithError(err).WithField("event", "error.msghandler.process_raw_notification.ProcessRawNotification.failed").Error()
				return nil, err
			}

			notificationReceived := subscription.NotificationReceived{
				NotificationUUID: rawNotification.UUID(),
			}

			notificationReceivedPayload, err := json.Marshal(notificationReceived)
			if err != nil {
				log.WithError(err).WithField("event", "error.msghandler.process_raw_notification.marshal-notification.failed").Error()
				return nil, err
			}

			notificationReceivedMsg := message.NewMessage(uuid.NewString(), notificationReceivedPayload)

			return message.Messages{notificationReceivedMsg}, nil

		},
	)

	router.AddHandler(
		"fetch_notification_data",
		topic.WithingsNotificationReceived,
		wrk.Subscriber,
		topic.WithingsNotificationDataFetched,
		wrk.Publisher,
		func(msg *message.Message) ([]*message.Message, error) {

			log = log.
				WithField("handler_name", "fetch_notification_data").
				WithField("message_uuid", msg.UUID).
				WithField("message_payload", string(msg.Payload)).
				WithField("message_metadata", msg.Metadata)

			log.WithField("event", "info.msghandler.fetch_notification_data.started").
				Info()

			var notificationReceived subscription.NotificationReceived
			err := json.Unmarshal(msg.Payload, &notificationReceived)
			if err != nil {
				log.WithError(err).WithField("event", "error.msghandler.fetch_notification_data.unmarshal.failed").Error()
				return nil, err
			}

			notification, err := wrk.SubscriptionRepo.GetNotificationByUUID(ctx, notificationReceived.NotificationUUID)
			if err != nil {
				log.WithError(err).WithField("event", "error.msghandler.fetch_notification_data.GetNotificationByUUID.failed").Error()
				return nil, err
			}

			err = wrk.Commands.FetchNotificationData.Handle(ctx, command.FetchNotificationData{
				Notification: notification,
			})
			if err != nil {
				log.WithError(err).WithField("event", "error.msghandler.fetch_notification_data.ProcessRawNotification.failed").Error()
				return nil, err
			}

			notificationDataFetched := subscription.NotificationDataFetched{
				NotificationUUID: notification.UUID(),
			}

			notificationDataFetchedPayload, err := json.Marshal(notificationDataFetched)
			if err != nil {
				log.WithError(err).WithField("event", "error.msghandler.fetch_notification_data.marshal-notification.failed").Error()
				return nil, err
			}

			notificationDataFetchedMsg := message.NewMessage(uuid.NewString(), notificationDataFetchedPayload)

			return message.Messages{notificationDataFetchedMsg}, nil

		},
	)

	if err := router.Run(ctx); err != nil {
		log.WithError(err).WithField("event", "panic.worker.router.Run.failed").Error()
		panic(err)
	}
}
