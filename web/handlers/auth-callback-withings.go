package handlers

import (
	"github.com/roessland/withoutings/internal/logging"
	"github.com/roessland/withoutings/internal/repos/db"
	"github.com/roessland/withoutings/internal/services/withoutings"
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

		// Validate state
		storedState := svc.Sessions.GetString(ctx, "state")
		callbackState := r.Form.Get("state")
		if callbackState != storedState || callbackState == "" {
			log.Infof("invalid state, had %s, expected %s", storedState, callbackState)
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
		svc.Sessions.Remove(ctx, "state")

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
		svc.Sessions.Put(ctx, "account_id", account.AccountID)

		// Redirect to homepage
		http.Redirect(w, r, "/", http.StatusFound)
	}
}
