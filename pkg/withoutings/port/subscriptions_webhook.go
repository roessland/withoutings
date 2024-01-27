package port

import (
	"encoding/json"
	"fmt"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/roessland/withoutings/pkg/logging"
	"github.com/roessland/withoutings/pkg/withoutings/adapter/topic"
	"github.com/roessland/withoutings/pkg/withoutings/app"
	"github.com/roessland/withoutings/pkg/withoutings/domain/subscription"
	"io"
	"net/http"
	"time"
)

// WithingsWebhook is the endpoint that Withings will send notifications to.
// There is no authentication, but the URL is secret, and IP whitelisting
// will be added later.
//
// Methods: HEAD, POST
func WithingsWebhook(svc *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		log := logging.MustGetLoggerFromContext(ctx)
		vars := mux.Vars(r)
		webhookSecret := vars["webhook_secret"]
		if webhookSecret != svc.Config.WithingsWebhookSecret {
			log.Error("invalid webhook URL secret")
			w.WriteHeader(401)
			return
		}

		data, err := io.ReadAll(r.Body)
		if err != nil {
			log.WithError(err).
				WithField("event", "error.withingswebhook.readbody.failed").
				Error()
			w.WriteHeader(500)
			return
		}

		// TODO add IP filtering
		// TODO log more headers
		// TODO use a command instead of directly calling the repo
		source := fmt.Sprintf("ip=%s", r.Header.Get("X-Forwarded-For"))
		rawNotification := subscription.NewRawNotification(
			uuid.New(),
			source,
			string(data),
			subscription.RawNotificationStatusPending,
			time.Now(),
			nil,
		)
		err = svc.SubscriptionRepo.CreateRawNotification(ctx,
			rawNotification,
		)
		if err != nil {
			log.WithError(err).
				WithField("event", "error.withingswebhook.createrawnotification.failed").
				Error()
			w.WriteHeader(500)
			return
		}

		// Emit event that notification was received
		rawNotificationReceived := subscription.RawNotificationReceived{
			RawNotificationUUID: rawNotification.UUID(),
		}
		msg, err := json.Marshal(rawNotificationReceived)
		if err != nil {
			log.WithError(err).
				WithField("event", "error.command.ProcessRawNotification.event.Marshal.failed").
				Error()
			w.WriteHeader(500)
		}
		log.Debug("Publishing event: ", string(msg))

		err = svc.Publisher.Publish(topic.WithingsRawNotificationReceived, message.NewMessage(uuid.NewString(), msg))
		if err != nil {
			log.WithError(err).
				WithField("event", "error.command.ProcessRawNotification.event.Publish.failed").
				Error()
			w.WriteHeader(500)
		}
		log.Debug("Published event: ", string(msg))

		w.WriteHeader(200)
	}
}
