package handlers

import (
	"github.com/roessland/withoutings/internal/logging"
	"github.com/roessland/withoutings/internal/service"
	"github.com/roessland/withoutings/web/middleware"
	"net/http"
)

func Homepage(svc *service.App) http.HandlerFunc {
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
