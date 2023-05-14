package port

import (
	"github.com/roessland/withoutings/pkg/withoutings/app"
	"net/http"
)

// NotificationsPage renders a list of all received notifications
// for the current account.
//
// Methods: GET
func NotificationsPage(svc *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		//ctx := r.Context()
		//log := logging.MustGetLoggerFromContext(ctx)
	}
}
