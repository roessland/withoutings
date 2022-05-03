package handlers

import (
	"bytes"
	"encoding/json"
	"github.com/roessland/withoutings/middleware"
	"github.com/roessland/withoutings/ptrof"
	"github.com/roessland/withoutings/server/app"
	"github.com/roessland/withoutings/withingsapi2/openapi2"
	"io"
	"io/ioutil"
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

		// Get sleep summary
		params := &openapi2.Sleepv2GetsummaryParams{
			Startdateymd:  "2021-09-01",
			Enddateymd:    "2021-10-01",
			Lastupdate:    0,
			DataFields:    ptrof.String("total_sleep_time"),
			Authorization: "Bearer " + token.AccessToken,
		}

		httpResp, err := app.WithingsClient.API2.Sleepv2Getsummary(ctx, params)
		buf, _ := ioutil.ReadAll(httpResp.Body)
		httpResp.Body = io.NopCloser(bytes.NewReader(buf))
		if err != nil {
			log.WithError(err).
				WithField("status", httpResp.StatusCode).
				WithField("body", string(buf)).
				Error("fetching sleep data: couldn't fetch")
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		err = openapi2.ParseErrorResponse(httpResp)
		if err != nil {
			log.WithError(err).
				Error("fetching sleep data: invalid status")
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		// Decode sleep summary into struct
		apiResp, err := openapi2.ParseSleepv2GetsummaryResponse(httpResp)
		if err != nil {
			log.WithError(err).Error("parsing sleep data response")
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		buf, err = json.Marshal(apiResp.JSON200)

		w.Header().Set("Content-Type", "text/html")
		err = app.Templates.RenderHomePage(w, token, string(buf))
		if err != nil {
			app.Log.WithError(err).WithField("event", "error.render.template").Error()
			return
		}
	}
}
