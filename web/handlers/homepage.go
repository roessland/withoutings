package handlers

import (
	"github.com/roessland/withoutings/app/webapp"
	"github.com/roessland/withoutings/logging"
	"net/http"
)

func HomePage(app *webapp.WebApp) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		log := logging.MustGetLoggerFromContext(ctx)

		sess, err := app.Sessions.Get(r)
		if err != nil {
			log.WithError(err).Error("parsing cookie")
			http.Error(w, "Invalid cookie", http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "text/html")
		err = app.Templates.RenderHomePage(w, sess.Token())
		if err != nil {
			app.Log.WithError(err).WithField("event", "error.render.template").Error()
			return
		}
	}
}
