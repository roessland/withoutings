package web_test

import (
	"github.com/gorilla/mux"
	"github.com/roessland/withoutings/pkg/withoutings/app"
	"github.com/roessland/withoutings/web"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestMiddlewareConfig(t *testing.T) {
	svc := app.NewMockApplication(t)
	router := mux.NewRouter()
	router.Use(web.Middleware(svc)...)

	var handler http.HandlerFunc

	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		handler(w, r)
	})

	t.Run("overwrites RemoteAddr with IP from X-Forwarded-For header", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.RemoteAddr = "192.168.1.1:1234"
		req.Header.Set("X-Forwarded-For", "123.123.123.123")
		handler = func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "123.123.123.123", r.RemoteAddr)
		}
		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)
	})
}
