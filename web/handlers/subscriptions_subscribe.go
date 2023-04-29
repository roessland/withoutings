package handlers

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/roessland/withoutings/pkg/logging"
	"github.com/roessland/withoutings/pkg/withoutings/app"
	"github.com/roessland/withoutings/pkg/withoutings/app/command"
	"github.com/roessland/withoutings/web/middleware"
	"net/http"
	"strconv"
)

// Subscribe subscribes to a single webhook category for the current user.
//
// Methods: POST
func Subscribe(svc *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		log := logging.MustGetLoggerFromContext(ctx).
			WithField("handler", "Subscribe")
		vars := mux.Vars(r)

		appli, err := strconv.Atoi(vars["appli"])
		if err != nil || appli == 0 {
			log.WithError(err).
				WithField("appli", vars["appli"]).
				WithField("event", "warn.illegal_parameter").
				Warn()
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "appli parameter must be an integer")
			return
		}

		account := middleware.GetAccountFromContext(ctx)
		if account == nil {
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprintf(w, "You must be logged in to subscribe to webhooks.")
			return
		}

		err = svc.Commands.SubscribeAccount.Handle(ctx, command.SubscribeAccount{
			Account: *account,
			Appli:   appli,
		})
		if err != nil {
			log.WithError(err).Error()
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "An error occurred when trying to subscribe to webhooks.")

			return
		}

		http.Redirect(w, r, "/subscriptions", http.StatusSeeOther)
	}
}
