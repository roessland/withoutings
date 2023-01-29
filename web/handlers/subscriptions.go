package handlers

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/roessland/withoutings/pkg/logging"
	"github.com/roessland/withoutings/pkg/withoutings/app"
	"github.com/roessland/withoutings/pkg/withoutings/app/command"
	"github.com/roessland/withoutings/pkg/withoutings/domain/subscription"
	"github.com/roessland/withoutings/web/middleware"
	"io"
	"net/http"
)

func SubscriptionsPage(svc *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		log := logging.MustGetLoggerFromContext(ctx)

		account := middleware.GetAccountFromContext(ctx)
		if account == nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		subscriptions, err := svc.SubscriptionRepo.GetSubscriptionsByAccountID(ctx, account.AccountID)
		w.Header().Set("Content-Type", "text/html")
		err = svc.Templates.RenderSubscriptionsPage(w, subscriptions)
		if err != nil {
			log.WithError(err).WithField("event", "error.render.template").Error()
			return
		}
	}
}

func Subscribe(svc *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		log := logging.MustGetLoggerFromContext(ctx)

		account := middleware.GetAccountFromContext(ctx)
		if account == nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		err := svc.Commands.SubscribeAccount.Handle(ctx, command.SubscribeAccount{
			Account: *account,
		})
		if err != nil {
			log.WithError(err).Error()
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/subscriptions", http.StatusSeeOther)
	}
}

func WithingsWebhook(svc *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		log := logging.MustGetLoggerFromContext(ctx)
		vars := mux.Vars(r)
		webhookSecret := vars["webhook_secret"]
		if webhookSecret != "supersecret" {
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
		source := fmt.Sprintf("ip=%s", r.RemoteAddr)
		err = svc.SubscriptionRepo.CreateRawNotification(ctx, subscription.NewRawNotification(source, string(data)))
		if err != nil {
			log.WithError(err).Error()
			w.WriteHeader(500)
			return
		}

		w.WriteHeader(200)
	}
}
