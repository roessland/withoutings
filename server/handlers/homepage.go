package handlers

import (
	"github.com/roessland/withoutings/middleware"
	"github.com/roessland/withoutings/server/app"
	"net/http"
)

func HomePage(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		log := middleware.MustGetLoggerFromContext(ctx)

		sess, err := app.Sessions.Get(r)
		if err != nil {
			log.WithError(err).Error("parsing cookie")
			http.Error(w, "Invalid cookie", http.StatusBadRequest)
			return
		}

		token := sess.Token()

		w.Header().Set("Content-Type", "text/html")
		err = app.Templates.RenderHomePage(w, token)
		if err != nil {
			app.Log.WithError(err).WithField("event", "error.render.template").Error()
			return
		}
	}
}
