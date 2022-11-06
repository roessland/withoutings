package handlers

import (
	"github.com/roessland/withoutings/app/webapp"
	"github.com/roessland/withoutings/logging"
	"net/http"
)

// Logout logs users out via Withings OAuth2.
func Logout(app *webapp.WebApp) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		log := logging.MustGetLoggerFromContext(ctx)

		sess, err := app.Sessions.Get(r)
		sess.Options.MaxAge = -1
		if err != nil {
			log.WithField("event", "error.logout.getsession").
				WithError(err).Error()
			http.Error(w, "Invalid cookie", http.StatusBadRequest)
			return
		}

		err = sess.Save(r, w)
		if err != nil {
			log.WithField("event", "error.logout.savesession").
				WithError(err).Error()
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/", http.StatusFound)
	}
}
