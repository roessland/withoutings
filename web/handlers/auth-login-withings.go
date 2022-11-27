package handlers

import (
	"github.com/roessland/withoutings/internal/service"
	"github.com/roessland/withoutings/internal/withoutings/adapters/withingsapi"
	"net/http"
)

// Login logs users in via Withings OAuth2.
func Login(svc *service.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		// log := logging.MustGetLoggerFromContext(ctx)

		// Save state to cookie. It will be verified in the callback handler.
		nonce := withingsapiadapter.RandomNonce()
		svc.Sessions.Put(ctx, "state", nonce)

		authCodeURL := svc.Withings.AuthCodeURL(nonce)
		http.Redirect(w, r, authCodeURL, http.StatusFound)
	}
}
