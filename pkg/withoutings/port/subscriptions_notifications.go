package port

import (
	"fmt"
	"github.com/roessland/withoutings/pkg/logging"
	"github.com/roessland/withoutings/pkg/withoutings/app"
	"github.com/roessland/withoutings/pkg/withoutings/domain/account"
	"net/http"
)

// NotificationsPage renders a list of all received notifications
// for the current account.
//
// Methods: GET
func NotificationsPage(svc *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		log := logging.MustGetLoggerFromContext(ctx)

		acc := account.GetFromContext(ctx)
		if acc == nil {
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprintf(w, "You must be logged in to show notifications.")
			return
		}
		log = log.WithField("account_uuid", acc.UUID())

		notifications, err := svc.SubscriptionRepo.GetNotificationsByAccountUUID(ctx, acc.UUID())
		if err != nil {
			log.WithField("event", "error.NotificationsPage.GetNotificationsByAccountUUID").
				WithError(err).
				Error()
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "An error occurred when trying to get notifications.")
			return
		}

		log.WithField("notifications", notifications).WithField("event", "debug.NotificationsPage.got-notifications").Debug()

		err = svc.Templates.RenderNotifications(ctx, w, notifications, "")
		if err != nil {
			log.WithField("event", "error.NotificationsPage.RenderNotifications").
				WithError(err).
				Error()
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "An error occurred when trying to render notifications.")
			return
		}
	}
}
