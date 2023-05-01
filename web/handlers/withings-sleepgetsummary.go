package handlers

import (
	"github.com/roessland/withoutings/pkg/logging"
	"github.com/roessland/withoutings/pkg/service/sleep"
	"github.com/roessland/withoutings/pkg/withoutings/app"
	"github.com/roessland/withoutings/web/middleware"
	"net/http"
	"time"
)

// SleepSummaries renders the sleep summaries page, showing the user's sleep for the last N days.
func SleepSummaries(svc *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		log := logging.MustGetLoggerFromContext(ctx)

		account := middleware.GetAccountFromContext(ctx)
		if account == nil {
			w.Header().Set("Content-Type", "text/html")
			w.WriteHeader(403)
			err := svc.Templates.RenderSleepSummaries(w, nil, "You must be logged in to view this page.")
			if err != nil {
				log.WithError(err).WithField("event", "error.render.template").Error()
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				return
			}
			return
		}

		var sleepData sleep.GetSleepSummaryOutput
		if time.Now().After(account.WithingsAccessTokenExpiry()) {
			w.Header().Set("Content-Type", "text/html")
			w.WriteHeader(200)
			err := svc.Templates.RenderSleepSummaries(w, nil, "Your token is expired. ")
			if err != nil {
				log.WithError(err).WithField("event", "error.render.template").Error()
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				return
			}
			return
		}

		sleepData, err := svc.Sleep.GetSleepSummaries(ctx, account.WithingsAccessToken(), sleep.GetSleepSummaryInput{
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
		err = svc.Templates.RenderSleepSummaries(w, &sleepData, "")
		if err != nil {
			log.WithError(err).WithField("event", "error.render.template").Error()
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
	}

}
