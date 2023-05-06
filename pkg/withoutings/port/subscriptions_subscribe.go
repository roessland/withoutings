package port

import (
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/roessland/withoutings/pkg/logging"
	"github.com/roessland/withoutings/pkg/withoutings/app"
	"github.com/roessland/withoutings/pkg/withoutings/app/command"
	"github.com/roessland/withoutings/pkg/withoutings/domain/account"
	"github.com/roessland/withoutings/pkg/withoutings/domain/subscription"
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

		acc := account.GetFromContext(ctx)
		if acc == nil {
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprintf(w, "You must be logged in to subscribe to webhooks.")
			return
		}

		err = svc.Commands.SubscribeAccount.Handle(ctx, command.SubscribeAccount{
			Account: *acc,
			Appli:   appli,
		})
		if errors.Is(err, subscription.ErrSubscriptionAlreadyExists) {
			svc.Flash.PutMsg(ctx, "You are already subscribed to this category.")
		} else if err != nil {
			log.WithError(err).Error()
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "An error occurred when trying to subscribe to webhooks.")
			return
		}

		http.Redirect(w, r, "/subscriptions", http.StatusSeeOther)
	}
}
