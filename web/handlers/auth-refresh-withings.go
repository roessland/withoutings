package handlers

import (
	"context"
	"errors"
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

		if !accInitial.CanRefreshAccessToken() {
			w.WriteHeader(200)
			tmplErr := svc.Templates.RenderRefreshAccessToken(w, nil,
				"Not refreshing your access token since it not yet expired.")
			if tmplErr != nil {
				log.WithError(tmplErr).WithField("event", "error.render.template").Error()
				return
			}
			return
		}

		newToken, err := svc.WithingsRepo.RefreshAccessToken(ctx, accInitial.WithingsRefreshToken)
		if err != nil {
			w.WriteHeader(500)
			tmplErr := svc.Templates.RenderRefreshAccessToken(w, newToken,
				"Could not refresh your access token since an error occurred.")
			if tmplErr != nil {
				log.WithError(tmplErr).WithField("event", "error.render.template").Error()
				return
			}
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
			tmplErr := svc.Templates.RenderRefreshAccessToken(w, newToken,
				"Could not update your access token since an error occurred.")
			if tmplErr != nil {
				log.WithError(tmplErr).WithField("event", "error.render.template").Error()
				return
			}
			return
		}

		w.Header().Set("Content-Type", "text/html")
		tmplErr := svc.Templates.RenderRefreshAccessToken(w, newToken, "")
		if tmplErr != nil {
			log.WithError(tmplErr).WithField("event", "error.render.template").Error()
			return
		}
	}
}
