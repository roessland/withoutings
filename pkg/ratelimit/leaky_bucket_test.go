package ratelimit_test

import (
	"context"
	"testing"
	"time"

	"github.com/roessland/withoutings/pkg/ratelimit"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/time/rate"
)

type Limiter interface {
	Allow() bool
	Wait(context.Context) error
}

func TestSlidingWindow(t *testing.T) {
	var rl Limiter

	t.Run("rate limits according to limit", func(t *testing.T) {
		rl = ratelimit.NewLeakyBucket(rate.Every(10 * time.Millisecond))

		ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
		defer cancel()

		count := 0

		// Take as many tokens as possible within the time limit
	free:
		for {
			select {
			case <-ctx.Done():
				break free
			default:
				// Block until a token is available
				rl.Wait(ctx)
				count++

				// Immediately after taking at token, it's not possible to take another.
				require.False(t, rl.Allow())
			}
		}

		assert.InDelta(t, 50, count, 15, "1/10ms * 500ms = 50")
	})
}
