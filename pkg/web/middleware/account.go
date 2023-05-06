package middleware

import (
	"errors"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/roessland/withoutings/pkg/logging"
	"github.com/roessland/withoutings/pkg/withoutings/app"
	"github.com/roessland/withoutings/pkg/withoutings/app/query"
	"github.com/roessland/withoutings/pkg/withoutings/domain/account"
	"net/http"
)

func Account(svc *app.App) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			log := logging.MustGetLoggerFromContext(ctx)

			accountUUID, _ := uuid.Parse(svc.Sessions.GetString(ctx, "account_uuid"))

			acc, err := svc.Queries.AccountByAccountUUID.Handle(ctx, query.AccountByUUID{
				AccountUUID: accountUUID,
			})
			if err != nil && !errors.Is(err, account.ErrAccountNotFound) {
				log.WithError(err).WithField("event", "error.getaccount").Error()
				http.Error(w, "Internal server error", http.StatusInternalServerError)
				return
			}
			if acc != nil && acc.UUID() != uuid.Nil {
				ctx = account.AddToContext(ctx, acc)
				ctx = logging.AddLoggerToContext(ctx, log.WithField("account_uuid", acc.UUID()))
			}
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
