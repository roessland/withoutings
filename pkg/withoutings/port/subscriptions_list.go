package port

import (
	"fmt"
	"github.com/roessland/withoutings/pkg/logging"
	"github.com/roessland/withoutings/pkg/withoutings/app"
	"github.com/roessland/withoutings/pkg/withoutings/domain/account"
	"net/http"
)

// SubscriptionsPage renders the subscriptions page, representing
// the database state of each subscription.
//
// Methods: GET
func SubscriptionsPage(svc *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		log := logging.MustGetLoggerFromContext(ctx)

		acc := account.GetFromContext(ctx)
		if acc == nil {
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
		subscriptions, err := svc.SubscriptionRepo.GetSubscriptionsByAccountUUID(ctx, acc.UUID())
		if err != nil {
			log.WithError(err).WithField("event", "error.AllNotificationCategories").Error()
			tmplErr := svc.Templates.RenderSubscriptionsPage(ctx, w, nil,
				categories, "An error occurred when retrieving your webhook subscriptions.")
			if tmplErr != nil {
				log.WithError(tmplErr).WithField("event", "error.RenderSubscriptionsPage").Error()
				return
			}
			return
		}

		w.WriteHeader(200)
		w.Header().Set("Content-Type", "text/html")
		tmplErr := svc.Templates.RenderSubscriptionsPage(ctx, w, subscriptions, categories, "")
		if tmplErr != nil {
			log.WithError(tmplErr).WithField("event", "error.RenderSubscriptionsPage").Error()
			return
		}
	}
}
