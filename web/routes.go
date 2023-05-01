package web

import (
	"github.com/gorilla/mux"
	"github.com/roessland/withoutings/pkg/withoutings/app"
	"github.com/roessland/withoutings/web/handlers"
	"net/http"
	"time"
)

func Router(svc *app.App) *mux.Router {
	svc.Validate()

	r := mux.NewRouter()
	r.HandleFunc("/api/health", handlers.Health(svc))
	r.HandleFunc("/", handlers.Homepage(svc))
	r.PathPrefix("/static/").Handler(handlers.Static(svc))

	r.Path("/auth/login").Methods("GET").Handler(handlers.Login(svc))
	r.HandleFunc("/auth/logout", handlers.Logout(svc))
	r.HandleFunc("/auth/callback", handlers.Callback(svc))
	r.HandleFunc("/auth/refresh", handlers.RefreshWithingsAccessToken(svc))

	r.HandleFunc("/sleepsummaries", handlers.SleepSummaries(svc))
	//r.HandleFunc("/sleepget.json", handlers.SleepGetJSON(svc))
	r.Path("/subscriptions").Methods("GET").Handler(handlers.SubscriptionsPage(svc))
	r.Path("/subscriptions/withings").Methods("GET").Handler(handlers.SubscriptionsWithingsPage(svc))

	r.Path("/subscriptions/subscribe/{appli}").Methods("POST").Handler(handlers.Subscribe(svc))
	r.Path("/withings/webhooks/{webhook_secret}").Methods("HEAD", "POST").Handler(handlers.WithingsWebhook(svc))

	r.Use(Middleware(svc)...)
	return r
}

func Server(svc *app.App) *http.Server {
	return &http.Server{
		Handler:      Router(svc),
		Addr:         svc.Config.ListenAddr,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
}
