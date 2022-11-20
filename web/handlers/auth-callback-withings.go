package handlers

import (
	"github.com/roessland/withoutings/internal/domain/services/withoutings"
	"github.com/roessland/withoutings/internal/logging"
	"github.com/roessland/withoutings/internal/repos/db"
	"net/http"
)

// Callback is used for OAuth2 callbacks,
// but also for event notifications.
func Callback(svc *withoutings.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		log := logging.MustGetLoggerFromContext(ctx)

		err := r.ParseForm()
		if err != nil {
			log.WithError(err).Error("parsing form")
			http.Error(w, "Error parsing form", http.StatusBadRequest)
			return
		}

		sess, err := svc.Sessions.Get(r)
		if err != nil {
			log.WithError(err).Error("parsing cookie")
			http.Error(w, "Invalid cookie", http.StatusBadRequest)
			return
		}

		// Validate state
		storedState := sess.State()
		state := r.Form.Get("state")
		if state != storedState || state == "" {
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
		token, err := svc.Withings.GetAccessToken(ctx, code)
		if err != nil {
			log.WithError(err).
				WithField("event", "error.callback.getaccesstoken").
				Info()
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Clear nonce
		sess.SetState("")

		// Create account
		account, err := svc.AccountRepo.CreateAccount(ctx, db.CreateAccountParams{
			WithingsUserID:            token.UserID,
			WithingsAccessToken:       token.AccessToken,
			WithingsRefreshToken:      token.RefreshToken,
			WithingsAccessTokenExpiry: token.Expiry,
			WithingsScopes:            token.Scope,
		})
		if err != nil {
			log.WithError(err).
				WithField("event", "error.callback.createaccount").
				Error()
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Login user by saving account ID to session.
		sess.SetAccountID(account.AccountID)

		// Save token // TODO remove
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
