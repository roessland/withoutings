package handlers

import (
	"fmt"
	"github.com/roessland/withoutings/pkg/logging"
	"github.com/roessland/withoutings/pkg/withoutings/app"
	"github.com/roessland/withoutings/pkg/withoutings/domain/withings"
	"github.com/roessland/withoutings/web/middleware"
	"net/http"
)

// SubscriptionsWithingsPage that queries the Withings API for each
// subscription category and displays which of them the user has subscribed to.
//
// Methods: GET
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
