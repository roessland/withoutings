package handlers

import (
	"github.com/google/uuid"
	"github.com/roessland/withoutings/pkg/logging"
	"github.com/roessland/withoutings/pkg/withoutings/app"
	"github.com/roessland/withoutings/pkg/withoutings/app/command"
	"github.com/roessland/withoutings/pkg/withoutings/app/query"
	"github.com/roessland/withoutings/pkg/withoutings/domain/account"
	"net/http"
)

// Callback is used for OAuth2 callbacks. It is called by Withings after the user has logged in.
//
// Methods: *
func Callback(svc *app.App) http.HandlerFunc {
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
		token, err := svc.WithingsRepo.GetAccessToken(ctx, code)
		if err != nil {
			log.WithError(err).
				WithField("event", "error.callback.getaccesstoken").
				Info()
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Clear nonce
		svc.Sessions.Remove(ctx, "state")

		// Create domain object with placeholder UUID
		acc, err := account.NewAccount(
			uuid.New(),
			token.UserID,
			token.AccessToken,
			token.RefreshToken,
			token.Expiry,
			token.Scope,
		)
		if err != nil {
			log.WithError(err).
				WithField("event", "error.callback.newaccount").
				Error()
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// If another account with same Withings user ID exists, replace it. UUID in DB will be preserved.
		// Otherwise create new account with UUID chosen above.
		err = svc.Commands.CreateOrUpdateAccount.Handle(ctx, command.CreateOrUpdateAccount{
			Account: acc,
		})
		if err != nil {
			log.WithError(err).
				WithField("event", "error.callback.createaccount").
				Error()
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Find account ID
		acc, err = svc.Queries.AccountByWithingsUserID.Handle(ctx, query.AccountByWithingsUserID{WithingsUserID: token.UserID})
		if err != nil {
			log.WithError(err).
				WithField("event", "error.callback.getaccount").
				Error()
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Login user by saving account ID to session.
		svc.Sessions.Put(ctx, "account_uuid", acc.UUID().String())

		// Redirect to homepage
		http.Redirect(w, r, "/", http.StatusFound)
	}
}
