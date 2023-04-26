package handlers

import (
	"github.com/roessland/withoutings/pkg/logging"
	"github.com/roessland/withoutings/pkg/service/sleep"
	"github.com/roessland/withoutings/pkg/withoutings/app"
	"github.com/roessland/withoutings/web/middleware"
	"net/http"
	"time"
)

func SleepSummaries(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		log := logging.MustGetLoggerFromContext(ctx)

		account := middleware.GetAccountFromContext(ctx)
		if account == nil {
			http.Error(w, "You must log in first", http.StatusUnauthorized)
			return
		}

		var sleepData sleep.GetSleepSummaryOutput
		if time.Now().After(account.WithingsAccessTokenExpiry()) {
			w.Header().Set("Content-Type", "text/html")
			w.WriteHeader(200)
			err := app.Templates.RenderSleepSummaries(w, nil, "Your token is expired. ")
			if err != nil {
				log.WithError(err).WithField("event", "error.render.template").Error()
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				return
			}
			return
		}

		sleepData, err := app.Sleep.GetSleepSummaries(ctx, account.WithingsAccessToken(), sleep.GetSleepSummaryInput{
			Year:  0,
			Month: 0,
		})
		if err != nil {
			log.Error(err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
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
