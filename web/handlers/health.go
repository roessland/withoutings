package handlers

import (
	"github.com/roessland/withoutings/internal/domain/services/withoutings"
	"net/http"
)

func Health(app *withoutings.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
	}
}
