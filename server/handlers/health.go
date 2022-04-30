package handlers

import (
	"github.com/roessland/withoutings/server/app"
	"net/http"
)

func Health(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
	}
}
