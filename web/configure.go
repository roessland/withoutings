package web

import (
	"github.com/gorilla/mux"
	"github.com/roessland/withoutings/internal/service"
	"github.com/roessland/withoutings/web/handlers"
	"github.com/roessland/withoutings/web/middleware"
	"github.com/roessland/withoutings/web/static"
	"net/http"
	"time"
)

func Router(svc *service.App) *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/api/health", handlers.Health(svc))
	r.HandleFunc("/", handlers.Homepage(svc))
	r.PathPrefix("/static/").Handler(http.FileServer(http.FS(static.FS)))

	r.Path("/auth/login").Methods("GET").Handler(handlers.Login(svc))
	r.HandleFunc("/auth/logout", handlers.Logout(svc))
	r.HandleFunc("/auth/callback", handlers.Callback(svc))
	//r.HandleFunc("/auth/refresh", handlers.Refresh(svc))
	//
	r.HandleFunc("/sleepsummaries", handlers.SleepSummaries(svc))
	//r.HandleFunc("/sleepget.json", handlers.SleepGetJSON(svc))

	r.Use(
		middleware.Logging(svc),
		svc.Sessions.LoadAndSave,
		middleware.Account(svc),
	)
	return r
}

func Server(svc *service.App) *http.Server {
	return &http.Server{
		Handler:      Router(svc),
		Addr:         "127.0.0.1:3628",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
}
