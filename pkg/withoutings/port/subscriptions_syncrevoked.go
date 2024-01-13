package port

import (
	"fmt"
	"github.com/roessland/withoutings/pkg/logging"
	"github.com/roessland/withoutings/pkg/withoutings/app"
	"github.com/roessland/withoutings/pkg/withoutings/app/command"
	"github.com/roessland/withoutings/pkg/withoutings/domain/account"
	"net/http"
)

func SyncRevokedSubscriptions(svc *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		log := logging.MustGetLoggerFromContext(ctx).
			WithField("handler", "SyncRevokedSubscriptions")

		acc := account.GetFromContext(ctx)
		if acc == nil {
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprintf(w, "You must be logged in to sync your subscriptions.")
			return
		}

		err := svc.Commands.SyncRevokedSubscriptions.Handle(ctx, command.SyncRevokedSubscriptions{
			Account: acc,
		})
		if err != nil {
			log.WithError(err).WithField("event", "error.syncrevoked.command.failed").Error()
			w.WriteHeader(500)
			fmt.Fprintf(w, "An error occurred when trying to sync your subscriptions.")
			return
		}

		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(200)
		fmt.Fprintf(w, "Subscriptions synced successfully.")
	}
}
