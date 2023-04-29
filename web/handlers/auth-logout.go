package handlers

import (
	"github.com/roessland/withoutings/pkg/logging"
	"github.com/roessland/withoutings/pkg/withoutings/app"
	"net/http"
)

// Logout logs users out via Withings OAuth2.
//
// Methods: *
func Logout(svc *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		log := logging.MustGetLoggerFromContext(ctx)

		err := svc.Sessions.Destroy(ctx)
		if err != nil {
			log.WithField("event", "error.logout.destroy_session").
				WithError(err).Error()
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/", http.StatusFound)
	}
}
