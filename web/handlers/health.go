package handlers

import (
	"github.com/roessland/withoutings/internal/service"
	"net/http"
)

func Health(svc *service.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
	}
}
