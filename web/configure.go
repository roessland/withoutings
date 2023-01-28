package web

import (
	"github.com/gorilla/mux"
	"github.com/roessland/withoutings/pkg/withoutings/app"
	"github.com/roessland/withoutings/web/handlers"
	"github.com/roessland/withoutings/web/middleware"
	"github.com/roessland/withoutings/web/static"
	"net/http"
	"time"
)

func Router(svc *app.App) *mux.Router {
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

	r.Path("/subscriptions").Methods("GET").Handler(handlers.SubscriptionsPage(svc))
	r.Path("/subscriptions/subscribe").Methods("POST").Handler(handlers.Subscribe(svc))

	r.Path("/withings/webhooks/{webhook_secret}").Methods("POST").Handler(handlers.WithingsWebhook(svc))

	r.Use(
		middleware.Logging(svc),
		svc.Sessions.LoadAndSave,
		middleware.Account(svc),
	)
	return r
}

func Server(svc *app.App) *http.Server {
	return &http.Server{
		Handler:      Router(svc),
		Addr:         "127.0.0.1:3628",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
}
