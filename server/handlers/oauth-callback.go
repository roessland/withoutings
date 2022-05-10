package handlers

import (
	"github.com/roessland/withoutings/logging"
	"github.com/roessland/withoutings/server/serverapp"
	"net/http"
)

// Callback is used for OAuth2 callbacks,
// but also for event notifications.
func Callback(app *serverapp.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		log := logging.MustGetLoggerFromContext(ctx)

		err := r.ParseForm()
		if err != nil {
			log.WithError(err).Error("parsing form")
			http.Error(w, "Error parsing form", http.StatusBadRequest)
			return
		}

		sess, err := app.Sessions.Get(r)
		if err != nil {
			log.WithError(err).Error("parsing cookie")
			http.Error(w, "Invalid cookie", http.StatusBadRequest)
			return
		}

		// Validate state
		storedState := sess.State()
		state := r.Form.Get("state")
		if state != storedState {
			log.Info("invalid state")
			http.Error(w, "State invalid", http.StatusBadRequest)
			return
		}

		// Get token
		code := r.Form.Get("code")
		if code == "" {
			log.Info("code not found")
			http.Error(w, "Code not found", http.StatusBadRequest)
			return
		}
		token, err := app.Withings.GetAccessToken(ctx, code)
		if err != nil {
			log.WithError(err).
				WithField("event", "error.callback.getaccesstoken").
				Info()
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Clear nonce
		sess.SetState("")

		// Save token
		sess.SetToken(token)

		// Save session
		err = sess.Save(r, w)
		if err != nil {
			log.WithField("event", "error.callback.setcookie").
				WithError(err).Error()
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		// Redirect to homepage
		http.Redirect(w, r, "/", http.StatusFound)
	}
}
