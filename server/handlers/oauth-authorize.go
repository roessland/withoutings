package handlers

import (
	"bytes"
	"fmt"
	"github.com/roessland/withoutings/middleware"
	"github.com/roessland/withoutings/server/app"
	"io"
	"io/ioutil"
	"net/http"
)

// Callback is used for OAuth2 callbacks,
// but also for event notifications.
func Callback(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		log := middleware.MustGetLoggerFromContext(ctx)

		buf, err := ioutil.ReadAll(r.Body)
		r.Body = io.NopCloser(bytes.NewReader(buf))

		err = r.ParseForm()
		if err != nil {
			log.WithError(err).Error("parsing form")
			http.Error(w, "Error parsing form", http.StatusBadRequest)
			return
		}

		state := r.Form.Get("state")
		if state != "xyfdsfdsz" {
			log.Info("invalid state")
			http.Error(w, "State invalid", http.StatusBadRequest)
			return
		}

		code := r.Form.Get("code")
		if code == "" {
			log.Info("code not found")
			http.Error(w, "Code not found", http.StatusBadRequest)
			return
		}

		token, err := app.WithingsClient.GetAccessToken(ctx, code)
		fmt.Println(token)
		if err != nil {
			log.WithError(err).
				WithField("event", "error.token.exchange").
				Info()
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		log.Println(token)
	}
}
