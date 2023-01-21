package handlers

import (
	"github.com/roessland/withoutings/pkg/logging"
	"github.com/roessland/withoutings/pkg/withoutings/app"
	"github.com/roessland/withoutings/web/middleware"
	"net/http"
)

func Homepage(svc *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		log := logging.MustGetLoggerFromContext(ctx)

		account := middleware.GetAccountFromContext(ctx)

		w.Header().Set("Content-Type", "text/html")
		err := svc.Templates.RenderHomePage(w, account)
		if err != nil {
			log.WithError(err).WithField("event", "error.render.template").Error()
			return
		}
	}
}
