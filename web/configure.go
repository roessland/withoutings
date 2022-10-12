package web

import (
	"github.com/gorilla/mux"
	"github.com/roessland/withoutings/middleware"
	"github.com/roessland/withoutings/web/handlers"
	"github.com/roessland/withoutings/web/static"
	"github.com/roessland/withoutings/web/webapp"
	"net/http"
	"time"
)

func Configure(app *app.App) *http.Server {
	r := mux.NewRouter()
	r.HandleFunc("/api/health", handlers.Health(app))
	r.HandleFunc("/", handlers.HomePage(app))
	r.PathPrefix("/static/").Handler(http.FileServer(http.FS(static.FS)))

	r.HandleFunc("/auth/login", handlers.Login(app))
	r.HandleFunc("/auth/logout", handlers.Logout(app))
	r.HandleFunc("/auth/callback", handlers.Callback(app))
	r.HandleFunc("/auth/refresh", handlers.Refresh(app))

	r.HandleFunc("/sleepsummaries", handlers.SleepSummaries(app))

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
