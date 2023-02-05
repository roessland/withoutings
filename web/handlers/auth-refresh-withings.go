package handlers

import (
	"github.com/roessland/withoutings/pkg/logging"
	"github.com/roessland/withoutings/pkg/withoutings/app"
	"github.com/roessland/withoutings/pkg/withoutings/app/command"
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

		err := svc.Commands.RefreshAccessToken.Handle(ctx, command.RefreshAccessToken{Account: *accInitial})
		if err != nil {
			log.WithError(err).
				WithField("event", "event.handlers.RefreshWithingsAccessToken.failed").
				Error()
			w.WriteHeader(500)
			tmplErr := svc.Templates.RenderRefreshAccessToken(w, nil,
				"Could not refresh your access token since an error occurred.")
			if tmplErr != nil {
				log.WithError(tmplErr).WithField("event", "error.render.template").Error()
				return
			}
			return
		}

		w.Header().Set("Content-Type", "text/html")
		tmplErr := svc.Templates.RenderRefreshAccessToken(w, nil, "")
		if tmplErr != nil {
			log.WithError(tmplErr).WithField("event", "error.render.template").Error()
			return
		}
	}
}
