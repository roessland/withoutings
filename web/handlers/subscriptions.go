package handlers

import (
	"errors"
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

		// Make sure subscription exists before storing notification to database.
		sub, err := svc.SubscriptionRepo.GetSubscriptionByWebhookSecret(ctx, webhookSecret)
		if errors.Is(err, subscription.NotFoundError{}) {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		if err != nil {
			log.WithError(err).Error()
			w.WriteHeader(500)
			return
		}

		data, err := io.ReadAll(r.Body)
		if err != nil {
			log.WithError(err).Error()
			w.WriteHeader(500)
			return
		}

		source := fmt.Sprintf("subscription_id:%d", sub.SubscriptionID)
		err = svc.SubscriptionRepo.CreateRawNotification(ctx, subscription.NewRawNotification(source, string(data)))
		if err != nil {
			log.WithError(err).Error()
			w.WriteHeader(500)
			return
		}

		w.WriteHeader(200)
	}
}
