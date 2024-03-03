package port

import (
	"github.com/roessland/withoutings/pkg/withoutings/app"
	"github.com/roessland/withoutings/pkg/withoutings/domain/withings"
	"net/http"
)

// RedirectToWithingsLogin logs users in via Withings OAuth2.
// After login, they are redirected to the Withings callback handler.
func RedirectToWithingsLogin(svc *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		// log := logging.MustGetLoggerFromContext(ctx)

		// Save state to cookie. It will be verified in the callback handler.
		nonce := withings.RandomNonce()
		svc.Sessions.Put(ctx, "state", nonce)

		authCodeURL := svc.WithingsRepo.AuthCodeURL(nonce)
		http.Redirect(w, r, authCodeURL, http.StatusFound)
	}
}
