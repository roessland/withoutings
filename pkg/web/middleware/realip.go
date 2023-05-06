package middleware

import (
	"github.com/gorilla/mux"
	"net/http"
	"strings"
)

// UseRemoteAddrFromXForwardedFor is a middleware that
// sets the real IP address of the client
// based on the X-Forwarded-For header from Caddy.
// This should only be used when running behind Caddy and other
// reverse proxies that set the X-Forwarded-For header.
func UseRemoteAddrFromXForwardedFor() mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Header.Get("X-Forwarded-For") != "" {
				// X-Forwarded-For is a comma-separated list of IPs.
				// Each reverse proxy the incoming IP to the back of the list.
				ips := strings.Split(r.Header.Get("X-Forwarded-For"), ",")
				ip := strings.TrimSpace(ips[0])
				r.RemoteAddr = ip
			}
			next.ServeHTTP(w, r)
		})
	}
}
