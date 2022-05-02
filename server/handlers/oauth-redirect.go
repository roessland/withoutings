package handlers

import (
	"github.com/roessland/withoutings/middleware"
	"github.com/roessland/withoutings/server/app"
	"github.com/roessland/withoutings/withingsOld"
	"net/http"
)

// Login logs users in via Withings OAuth2.
func Login(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		log := middleware.MustGetLoggerFromContext(ctx)

		sess, err := app.Sessions.Get(r)
		if err != nil {
			log.WithField("event", "error.login.getsession").
				WithError(err).Error()
			http.Error(w, "Invalid cookie", http.StatusBadRequest)
			return
		}

		// Save state to cookie. It will be verified in the callback handler.
		nonce := withingsOld.RandomNonce()
		sess.SetState(nonce)
		err = sess.Save(r, w)
		if err != nil {
			log.WithField("event", "error.login.setcookie").
				WithError(err).Error()
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		authCodeURL := app.WithingsClient.AuthCodeURL(nonce)
		http.Redirect(w, r, authCodeURL, http.StatusFound)
	}
}
