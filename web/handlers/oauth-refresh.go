package handlers

import (
	"github.com/roessland/withoutings/domain/services/withoutings"
	"github.com/roessland/withoutings/logging"
	"net/http"
)

// Refresh gets a new OAuth2 token using refresh token.
func Refresh(app *withoutings.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		log := logging.MustGetLoggerFromContext(ctx)

		sess, err := app.Sessions.Get(r)
		if err != nil {
			log.WithError(err).Error("parsing cookie")
			http.Error(w, "Invalid cookie", http.StatusBadRequest)
			return
		}

		// Validate state
		token := sess.Token()
		if token == nil {
			http.Error(w, "You don't have a token to refresh", http.StatusBadRequest)
			return
		}

		prevToken := token

		// Refresh access token
		token, err = app.Withings.RefreshAccessToken(ctx, token.RefreshToken)
		if err != nil {
			log.WithError(err).
				WithField("event", "error.refreshaccesstoken").
				Info()
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Save token
		sess.SetToken(token)

		// Save session
		err = sess.Save(r, w)
		if err != nil {
			log.WithField("event", "error.refreshaccesstoken.setcookie").
				WithError(err).Error()
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		// Show refreshed message.
		w.Header().Set("Content-Type", "text/html")
		err = app.Templates.RenderRefreshAccessToken(w, token, prevToken)
		if err != nil {
			log.WithField("event", "error.refreshaccesstoken.render").
				WithError(err).Error()
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
	}
}
