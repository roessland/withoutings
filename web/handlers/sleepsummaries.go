package handlers

import (
	"github.com/roessland/withoutings/app/webapp"
	"github.com/roessland/withoutings/domain/services/sleep"
	"github.com/roessland/withoutings/logging"
	"net/http"
	"time"
)

func SleepSummaries(app *webapp.WebApp) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		log := logging.MustGetLoggerFromContext(ctx)

		sess, err := app.Sessions.Get(r)
		if err != nil {
			log.WithError(err).Error("parsing cookie")
			http.Error(w, "Invalid cookie", http.StatusBadRequest)
			return
		}

		var sleepData sleep.GetSleepSummaryOutput
		token := sess.Token()

		if token != nil {
			if time.Now().After(token.Expiry) {
				w.Header().Set("Content-Type", "text/html")
				w.WriteHeader(200)
				err = app.Templates.RenderSleepSummaries(w, nil, "Your token is expired. Go refresh it.")
				if err != nil {
					log.WithError(err).WithField("event", "error.render.template").Error()
					http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
					return
				}
				return
			}

			sleepData, err = app.Sleep.GetSleepSummaries(ctx, sleep.GetSleepSummaryInput{
				AccessToken: token.AccessToken,
				Year:        0,
				Month:       0,
			})
			if err != nil {
				log.Error(err)
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				return
			}
		}

		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(200)
		err = app.Templates.RenderSleepSummaries(w, &sleepData, "")
		if err != nil {
			log.WithError(err).WithField("event", "error.render.template").Error()
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
	}
}
