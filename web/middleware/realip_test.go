package middleware_test

import (
	"github.com/roessland/withoutings/web/middleware"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCaddyRealIP_XForwardedMissing_UsesRemoteAddr(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.RemoteAddr = "222.222.222.222"
	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "222.222.222.222", r.RemoteAddr)
		w.WriteHeader(http.StatusOK)
	})

	middleware.UseRemoteAddrFromXForwardedFor()(handler).ServeHTTP(rr, req)
	require.Equal(t, http.StatusOK, rr.Code)
}

func TestCaddyRealIP_XForwardedSingle_UsesHeader(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("X-Forwarded-For", "123.123.123.123")
	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "123.123.123.123", r.RemoteAddr)
		w.WriteHeader(http.StatusOK)
	})

	middleware.UseRemoteAddrFromXForwardedFor()(handler).ServeHTTP(rr, req)
	require.Equal(t, http.StatusOK, rr.Code)
}

func TestCaddyRealIP_XForwardedMultiple_UsesFirstIPInHeader(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("X-Forwarded-For", "123.123.123.123, 321.321.321.321")
	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "123.123.123.123", r.RemoteAddr)
		w.WriteHeader(http.StatusOK)
	})

	middleware.UseRemoteAddrFromXForwardedFor()(handler).ServeHTTP(rr, req)
	require.Equal(t, http.StatusOK, rr.Code)
}

func TestCaddyRealIP_XForwardedMultiple_UsesFirstHeader_FirstIP(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header["X-Forwarded-For"] = []string{
		"123.123.123.123, 321.321.321.321",
		"012.012.012.012, 210.210.210.210",
	}
	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "123.123.123.123", r.RemoteAddr)
		w.WriteHeader(http.StatusOK)
	})

	middleware.UseRemoteAddrFromXForwardedFor()(handler).ServeHTTP(rr, req)
	require.Equal(t, http.StatusOK, rr.Code)
}
