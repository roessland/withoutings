package handlers

import (
	"github.com/roessland/withoutings/pkg/withoutings/app"
	"net/http"
)

// Health is a simple health check endpoint.
//
// Methods: HEAD, GET
func Health(svc *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
	}
}
