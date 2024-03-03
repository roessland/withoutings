package withings

import "context"

type RateLimiter interface {
	Wait(ctx context.Context) error
}
