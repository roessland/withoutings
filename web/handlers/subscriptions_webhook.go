package handlers

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/roessland/withoutings/pkg/logging"
	"github.com/roessland/withoutings/pkg/withoutings/app"
	"github.com/roessland/withoutings/pkg/withoutings/domain/subscription"
	"io"
	"net/http"
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
			log.WithError(err).Error()
			w.WriteHeader(500)
			return
		}

		// TODO add IP filtering
		// TODO log more headers
		source := fmt.Sprintf("ip=%s", r.Header.Get("X-Forwarded-For"))
		err = svc.SubscriptionRepo.CreateRawNotification(ctx, subscription.NewRawNotification(source, string(data)))
		if err != nil {
			log.WithError(err).Error()
			w.WriteHeader(500)
			return
		}

		w.WriteHeader(200)
	}
}
