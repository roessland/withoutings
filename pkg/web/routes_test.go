package web_test

import (
	"github.com/roessland/withoutings/pkg/web"
	"github.com/roessland/withoutings/pkg/withoutings/app"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRouterConfig_DoesntCrash(t *testing.T) {
	svc := app.NewMockApplication(t)

	assertReturns200 := func(target string) {
		t.Run("GET "+target+" returns 200", func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, target, nil)
			rr := httptest.NewRecorder()
			router := web.Router(svc)
			router.ServeHTTP(rr, req)
			assert.EqualValues(t, http.StatusOK, rr.Code)
		})
	}

	assertReturns200("/api/health")
	assertReturns200("/favicon.ico")
	assertReturns200("/static/icon-512.png")
}
