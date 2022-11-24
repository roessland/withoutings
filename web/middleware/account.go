package middleware

import (
	"context"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v4"
	"github.com/roessland/withoutings/internal/logging"
	"github.com/roessland/withoutings/internal/repos/db"
	"github.com/roessland/withoutings/internal/service"
	"net/http"
)

var contextKeyAccount contextKey = "requestID"

func GetAccountFromContext(ctx context.Context) *db.Account {
	account, ok := ctx.Value(contextKeyAccount).(db.Account)
	if !ok {
		return nil
	}
	return &account
}

func AddAccountToContext(ctx context.Context, account db.Account) context.Context {
	return context.WithValue(ctx, contextKeyAccount, account)
}

func Account(svc *service.App) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			log := logging.MustGetLoggerFromContext(ctx)

			accountID := svc.Sessions.GetInt64(ctx, "account_id")
			account, err := svc.AccountRepo.GetAccount(ctx, accountID)
			if err != nil && err != pgx.ErrNoRows {
				log.WithError(err).WithField("event", "error.getaccount").Error()
				http.Error(w, "Internal server error", http.StatusInternalServerError)
				return
			}
			if account.AccountID != 0 {
				ctx = AddAccountToContext(ctx, account)
			}
			ctx = logging.AddLoggerToContext(ctx, log.WithField("account_id", account.AccountID))
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
