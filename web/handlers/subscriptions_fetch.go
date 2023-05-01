package handlers

import (
	"fmt"
	"github.com/roessland/withoutings/pkg/logging"
	"github.com/roessland/withoutings/pkg/withoutings/app"
	"github.com/roessland/withoutings/pkg/withoutings/domain/account"
	"github.com/roessland/withoutings/pkg/withoutings/domain/withings"
	"github.com/roessland/withoutings/web/templates"
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

		acc := account.GetAccountFromContext(ctx)
		if acc == nil {
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

		withingsSubscriptions := make([]templates.SubscriptionsWithingsPageItem, 0)
		for _, cat := range categories {
			notifyListResponse, err := svc.WithingsRepo.NotifyList(ctx, acc.WithingsAccessToken(),
				withings.NewNotifyListParams(cat.Appli))
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				fmt.Fprintf(w, "Error checking notification status with Withings")
				return
			}
			if len(notifyListResponse.Body.Profiles) == 0 {
				withingsSubscriptions = append(withingsSubscriptions, templates.SubscriptionsWithingsPageItem{
					Appli:            cat.Appli,
					AppliDescription: cat.Description,
					Exists:           false,
				})
			}
			for _, profile := range notifyListResponse.Body.Profiles {
				withingsSubscriptions = append(withingsSubscriptions, templates.SubscriptionsWithingsPageItem{
					Appli:            profile.Appli,
					AppliDescription: cat.Description,
					Exists:           true,
					Comment:          profile.Comment,
				})
			}
		}

		w.WriteHeader(200)
		w.Header().Set("Content-Type", "text/html")
		tmplErr := svc.Templates.RenderSubscriptionsWithingsPage(ctx, w, withingsSubscriptions, "")
		if tmplErr != nil {
			log.WithError(tmplErr).WithField("event", "error.RenderSubscriptionsWithingsPage").Error()
			return
		}
	}
}
