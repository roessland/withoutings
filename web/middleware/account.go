package middleware

import (
	"context"
	"errors"
	"github.com/gorilla/mux"
	"github.com/roessland/withoutings/pkg/logging"
	"github.com/roessland/withoutings/pkg/withoutings/app"
	"github.com/roessland/withoutings/pkg/withoutings/app/query"
	"github.com/roessland/withoutings/pkg/withoutings/domain/account"
	"net/http"
)

var contextKeyAccount contextKey = "requestID"

func GetAccountFromContext(ctx context.Context) *account.Account {
	acc, ok := ctx.Value(contextKeyAccount).(account.Account)
	if !ok {
		return nil
	}
	return &acc
}

func AddAccountToContext(ctx context.Context, account account.Account) context.Context {
	return context.WithValue(ctx, contextKeyAccount, account)
}

func Account(svc *app.App) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			log := logging.MustGetLoggerFromContext(ctx)

			accountID := svc.Sessions.GetInt64(ctx, "account_id")
			acc, err := svc.Queries.AccountForUserID.Handle(ctx, query.AccountByID{
				AccountID: accountID,
			})
			if err != nil && !errors.Is(err, account.NotFoundError{}) {
				log.WithError(err).WithField("event", "error.getaccount").Error()
				http.Error(w, "Internal server error", http.StatusInternalServerError)
				return
			}
			if acc.AccountID != 0 {
				ctx = AddAccountToContext(ctx, acc)
			}
			ctx = logging.AddLoggerToContext(ctx, log.WithField("account_id", acc.AccountID))
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
