package handlers

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/roessland/withoutings/pkg/logging"
	"github.com/roessland/withoutings/pkg/withoutings/app"
	"github.com/roessland/withoutings/pkg/withoutings/app/command"
	"github.com/roessland/withoutings/pkg/withoutings/domain/subscription"
	"github.com/roessland/withoutings/pkg/withoutings/domain/withings"
	"github.com/roessland/withoutings/web/middleware"
	"io"
	"net/http"
	"strconv"
)

func SubscriptionsPage(svc *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		log := logging.MustGetLoggerFromContext(ctx)

		account := middleware.GetAccountFromContext(ctx)
		if account == nil {
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprintf(w, "You must log in to show subscriptions")
			return
		}

		// List notification categories.
		categories, err := svc.SubscriptionRepo.AllNotificationCategories(ctx)
		if err != nil {
			log.WithError(err).WithField("event", "error.AllNotificationCategories").Error()
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "Error fetching notification categories")
			return
		}

		// Get persisted subscriptions.
		subscriptions, err := svc.SubscriptionRepo.GetSubscriptionsByAccountUUID(ctx, account.UUID())
		if err != nil {
			log.WithError(err).WithField("event", "error.AllNotificationCategories").Error()
			tmplErr := svc.Templates.RenderSubscriptionsPage(w, nil,
				categories, "An error occurred when retrieving your webhook subscriptions.")
			if tmplErr != nil {
				log.WithError(tmplErr).WithField("event", "error.RenderSubscriptionsPage").Error()
				return
			}
			return
		}

		w.WriteHeader(200)
		w.Header().Set("Content-Type", "text/html")
		tmplErr := svc.Templates.RenderSubscriptionsPage(w, subscriptions, categories, "")
		if tmplErr != nil {
			log.WithError(tmplErr).WithField("event", "error.RenderSubscriptionsPage").Error()
			return
		}
	}
}

func SubscriptionsWithingsPage(svc *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		log := logging.MustGetLoggerFromContext(ctx)

		account := middleware.GetAccountFromContext(ctx)
		if account == nil {
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprintf(w, "You must log in to show subscriptions")
			return
		}

		// Check WithingsRepo for each notification category.
		categories, err := svc.SubscriptionRepo.AllNotificationCategories(ctx)
		if err != nil {
			log.WithError(err).WithField("event", "error.AllNotificationCategories").Error()
			return
		}

		var withingsResponses []string
		for _, cat := range categories {
			notifyListResponse, err := svc.WithingsRepo.NotifyList(ctx, account.WithingsAccessToken(),
				withings.NewNotifyListParams(cat.Appli))
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				fmt.Fprintf(w, "Error checking notification status with Withings")
				return
			}
			withingsResponses = append(withingsResponses, string(notifyListResponse.Raw))
		}

		w.WriteHeader(200)
		w.Header().Set("Content-Type", "text/html")
		tmplErr := svc.Templates.RenderSubscriptionsWithingsPage(w, withingsResponses, "")
		if tmplErr != nil {
			log.WithError(tmplErr).WithField("event", "error.RenderSubscriptionsWithingsPage").Error()
			return
		}
	}
}

func Subscribe(svc *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		log := logging.MustGetLoggerFromContext(ctx).
			WithField("handler", "Subscribe")
		vars := mux.Vars(r)

		appli, err := strconv.Atoi(vars["appli"])
		if err != nil || appli == 0 {
			log.WithError(err).
				WithField("appli", vars["appli"]).
				WithField("event", "warn.illegal_parameter").
				Warn()
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "appli parameter must be an integer")
			return
		}

		account := middleware.GetAccountFromContext(ctx)
		if account == nil {
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprintf(w, "You must be logged in to subscribe to webhooks.")
			return
		}

		err = svc.Commands.SubscribeAccount.Handle(ctx, command.SubscribeAccount{
			Account: *account,
			Appli:   appli,
		})
		if err != nil {
			log.WithError(err).Error()
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "An error occurred when trying to subscribe to webhooks.")

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
		// TODO handle when we are behind proxy so we don't just log 127.0.0.1
		// TODO log more headers
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
