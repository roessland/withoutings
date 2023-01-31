package handlers

import (
	"context"
	"errors"
	"fmt"
	"github.com/roessland/withoutings/pkg/logging"
	"github.com/roessland/withoutings/pkg/withoutings/app"
	"github.com/roessland/withoutings/pkg/withoutings/domain/account"
	"github.com/roessland/withoutings/web/middleware"
	"net/http"
)

func RefreshWithingsAccessToken(svc *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		log := logging.MustGetLoggerFromContext(ctx)

		accInitial := middleware.GetAccountFromContext(ctx)
		if accInitial == nil {
			http.Error(w, "You must log in first", http.StatusUnauthorized)
			return
		}

		// TODO don't refresh if it's not expired yet.

		newToken, err := svc.Withings.RefreshAccessToken(ctx, accInitial.WithingsRefreshToken)
		if err != nil {
			w.WriteHeader(500)
			fmt.Fprintf(w, "couldn't refresh access token")
			return
		}

		err = svc.AccountRepo.UpdateAccount(
			ctx,
			accInitial.AccountID,
			func(ctx context.Context, accNext account.Account) (account.Account, error) {
				if accNext.WithingsRefreshToken != accInitial.WithingsRefreshToken {
					return account.Account{}, errors.New("refresh token updated by someone else")
				}
				accNext.WithingsAccessToken = newToken.AccessToken
				accNext.WithingsRefreshToken = newToken.RefreshToken
				accNext.WithingsAccessTokenExpiry = newToken.Expiry
				return accNext, nil
			},
		)
		if err != nil {
			w.WriteHeader(500)
			fmt.Fprintf(w, "couldn't refresh access token")
			return
		}

		w.Header().Set("Content-Type", "text/html")
		err = svc.Templates.RenderRefreshAccessToken(w, newToken)
		if err != nil {
			log.WithError(err).WithField("event", "error.render.template").Error()
			return
		}
	}
}
