package web_test

import (
	"github.com/roessland/withoutings/pkg/withoutings/app"
	"github.com/roessland/withoutings/web"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRouterConfig_DoesntCrash(t *testing.T) {
	svc := app.NewMockApplication(t)
	req := httptest.NewRequest(http.MethodGet, "/api/health", nil)
	rr := httptest.NewRecorder()
	router := web.Router(svc)
	router.ServeHTTP(rr, req)
}
