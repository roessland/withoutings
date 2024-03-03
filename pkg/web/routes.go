package web

import (
	"github.com/gorilla/mux"
	"github.com/roessland/withoutings/pkg/withoutings/app"
	"github.com/roessland/withoutings/pkg/withoutings/port"
	"net/http"
	"time"
)

func Router(svc *app.App) *mux.Router {
	svc.Validate()

	r := mux.NewRouter()
	r.HandleFunc("/api/health", port.Health(svc))
	r.HandleFunc("/", port.Homepage(svc))
	r.PathPrefix("/static/").Handler(port.Static(svc))
	r.Path("/favicon.ico").Handler(port.Static(svc))

	//r.Path("/login").Methods("GET").Handler(port.

	r.Path("/auth/login").Methods("GET").Handler(port.LoginPage(svc))
	r.Path("/auth/redirect-to-withings-login").Methods("POST").Handler(port.RedirectToWithingsLogin(svc))
	r.HandleFunc("/auth/logout", port.Logout(svc))
	r.HandleFunc("/auth/callback", port.WithingsCallback(svc))
	r.HandleFunc("/auth/refresh", port.RefreshWithingsAccessToken(svc))

	r.HandleFunc("/sleepsummaries", port.SleepSummaries(svc))
	//r.HandleFunc("/sleepget.json", handlers.SleepGetJSON(svc))
	r.Path("/subscriptions").Methods("GET").Handler(port.SubscriptionsPage(svc))
	r.Path("/subscriptions/withings").Methods("GET").Handler(port.SubscriptionsWithingsPage(svc))
	r.Path("/notifications").Methods("GET").Handler(port.NotificationsPage(svc))

	// TODO make it POST
	r.Path("/commands/sync-revoked-subscriptions").Methods("GET").Handler(port.SyncRevokedSubscriptions(svc))

	r.Path("/subscriptions/subscribe/{appli}").Methods("POST").Handler(port.Subscribe(svc))
	r.Path("/withings/webhooks/{webhook_secret}").Methods("HEAD", "POST").Handler(port.WithingsWebhook(svc))
	r.Path("/withings/measure/getmeas").Methods("POST").Handler(port.MeasureGetmeas(svc))

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
