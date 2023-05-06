package port

import (
	"github.com/roessland/withoutings/pkg/logging"
	"github.com/roessland/withoutings/pkg/withoutings/app"
	"github.com/roessland/withoutings/pkg/withoutings/domain/account"
	"net/http"
)

// Homepage renders the homepage.
//
// Methods: *
func Homepage(svc *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		log := logging.MustGetLoggerFromContext(ctx)

		acc := account.GetFromContext(ctx)

		w.Header().Set("Content-Type", "text/html")
		err := svc.Templates.RenderHomePage(ctx, w, acc)
		if err != nil {
			log.WithError(err).WithField("event", "error.render.template").Error()
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
	}
}
