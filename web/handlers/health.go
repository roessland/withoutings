package handlers

import (
	"github.com/roessland/withoutings/app/webapp"
	"net/http"
)

func Health(app *webapp.WebApp) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
	}
}
