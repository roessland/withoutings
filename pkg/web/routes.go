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

	// root
	// 	- globalMw(/static/)
	// 	- authMw(globalMw(/))

	r := mux.NewRouter()
	r.Use(GlobalMiddlewares(svc)...)

	r.PathPrefix("/static/").Handler(port.Static(svc))
	r.HandleFunc("/api/health", port.Health(svc))
	r.Path("/favicon.ico").Handler(port.Static(svc))

	s := r.PathPrefix("/").Subrouter()
	s.Use(AuthMiddlewares(svc)...)

	s.HandleFunc("/", port.Homepage(svc))

	s.Path("/auth/login").Methods("GET").Handler(port.LoginPage(svc))
	s.Path("/auth/redirect-to-withings-login").Methods("POST").Handler(port.RedirectToWithingsLogin(svc))
	s.HandleFunc("/auth/logout", port.Logout(svc))
	s.HandleFunc("/auth/callback", port.WithingsCallback(svc))
	s.HandleFunc("/auth/refresh", port.RefreshWithingsAccessToken(svc))

	s.HandleFunc("/sleepsummaries", port.SleepSummaries(svc))
	//r.HandleFunc("/sleepget.json", handlers.SleepGetJSON(svc))
	s.Path("/subscriptions").Methods("GET").Handler(port.SubscriptionsPage(svc))
	s.Path("/subscriptions/withings").Methods("GET").Handler(port.SubscriptionsWithingsPage(svc))
	s.Path("/notifications").Methods("GET").Handler(port.NotificationsPage(svc))

	// TODO make it POST
	s.Path("/commands/sync-revoked-subscriptions").Methods("GET").Handler(port.SyncRevokedSubscriptions(svc))

	s.Path("/subscriptions/subscribe/{appli}").Methods("POST").Handler(port.Subscribe(svc))
	s.Path("/withings/webhooks/{webhook_secret}").Methods("HEAD", "POST").Handler(port.WithingsWebhook(svc))
	s.Path("/withings/measure/getmeas").Methods("POST").Handler(port.MeasureGetmeas(svc))

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
