package port

import (
	"github.com/roessland/withoutings/pkg/logging"
	"github.com/roessland/withoutings/pkg/withoutings/app"
	"github.com/roessland/withoutings/pkg/withoutings/app/command"
	"github.com/roessland/withoutings/pkg/withoutings/domain/account"
	"net/http"
)

// RefreshWithingsAccessToken refreshes the access token for the current user.
// TODO: Before refreshing, account must be marked as "refreshing" in the database,
// so that we don't lose access to the account if the refresh fails, since
// the current access token will be invalidated after some time.
// A batch job must keep retrying to refresh the access token until it succeeds.
func RefreshWithingsAccessToken(svc *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		log := logging.MustGetLoggerFromContext(ctx)

		maybeAcc := account.GetFromContext(ctx)
		if maybeAcc == nil {
			http.Error(w, "You must log in first", http.StatusUnauthorized)
			return
		}

		if !maybeAcc.CanRefreshAccessToken() {
			w.WriteHeader(200)
			tmplErr := svc.Templates.RenderRefreshAccessToken(ctx, w, nil,
				"Not refreshing your access token since it not yet expired.")
			if tmplErr != nil {
				log.WithError(tmplErr).WithField("event", "error.render.template").Error()
				return
			}
			return
		}

		err := svc.Commands.RefreshAccessToken.Handle(ctx, command.RefreshAccessToken{Account: maybeAcc})
		if err != nil {
			log.WithError(err).
				WithField("event", "event.handlers.RefreshWithingsAccessToken.failed").
				Error()
			w.WriteHeader(500)
			tmplErr := svc.Templates.RenderRefreshAccessToken(ctx, w, nil,
				"Could not refresh your access token since an error occurred.")
			if tmplErr != nil {
				log.WithError(tmplErr).WithField("event", "error.render.template").Error()
				return
			}
			return
		}

		w.Header().Set("Content-Type", "text/html")
		tmplErr := svc.Templates.RenderRefreshAccessToken(ctx, w, nil, "")
		if tmplErr != nil {
			log.WithError(tmplErr).WithField("event", "error.render.template").Error()
			return
		}
	}
}
