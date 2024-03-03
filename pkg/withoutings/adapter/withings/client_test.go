package withings_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/roessland/withoutings/pkg/logging"
	withingsAdapter "github.com/roessland/withoutings/pkg/withoutings/adapter/withings"
	"github.com/roessland/withoutings/pkg/withoutings/domain/withings"
	"github.com/stretchr/testify/require"
)

func TestHTTPClient(t *testing.T) {
	ctx := context.Background()
	ctx = logging.AddLoggerToContext(ctx, logging.GetOrCreateLoggerFromContext(ctx))

	t.Run("Rate limits at fixed rate, using global rate limiter", func(t *testing.T) {
		var requestsReceived atomic.Int32
		mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestsReceived.Add(1)
			w.WriteHeader(http.StatusOK)
		}))

		ctx, cancel := context.WithTimeout(ctx, 1700*time.Millisecond)
		defer cancel()

		// Do as many requests as possible within the time limit, using two parallel clients.
		var wg sync.WaitGroup
		wg.Add(20)
		for i := 0; i < 20; i++ {
			go func() {
				defer wg.Done()
				c := withingsAdapter.NewClient("some-client-id", "some-client-secret", "some-redirect-url")
				c.APIBase = mockServer.URL
				for {
					select {
					case <-ctx.Done():
						return
					default:
						r, _ := c.NewNotifySubscribeRequest(withings.NewNotifySubscribeParams())
						r = r.WithContext(ctx)
						c.Do(r)
					}
				}
			}()
		}
		wg.Wait()

		// Using default rate limit of 120 requests per minute, we can do API
		// requests at a rate of 2 per second. As the bucket starts with 1 token,
		// we can instantly start request #1. #2, #3 and #4 then start at ~500 ms,
		// ~1000 ms, ~1500 ms and so on.
		require.InDelta(t, 4, requestsReceived.Load(), 1, "120/min * 700 ms = 2.8")
	})
}
