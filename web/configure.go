package web

import (
	"github.com/gorilla/mux"
	"github.com/roessland/withoutings/internal/domain/services/withoutings"
	"github.com/roessland/withoutings/web/handlers"
	"github.com/roessland/withoutings/web/middleware"
	"github.com/roessland/withoutings/web/static"
	"net/http"
	"time"
)

func Configure(app *withoutings.Service) *http.Server {
	r := mux.NewRouter()
	r.HandleFunc("/api/health", handlers.Health(app))
	r.HandleFunc("/", handlers.Homepage(app))
	r.PathPrefix("/static/").Handler(http.FileServer(http.FS(static.FS)))

	r.Path("/auth/login").Methods("GET").Handler(handlers.Login(app))
	r.HandleFunc("/auth/logout", handlers.Logout(app))
	r.HandleFunc("/auth/callback", handlers.Callback(app))
	r.HandleFunc("/auth/refresh", handlers.Refresh(app))

	r.HandleFunc("/sleepsummaries", handlers.SleepSummaries(app))
	r.HandleFunc("/sleepget.json", handlers.SleepGetJSON(app))

	r.Use(
		middleware.Logging(app),
	)

	return &http.Server{
		Handler:      r,
		Addr:         "127.0.0.1:3628",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
}
