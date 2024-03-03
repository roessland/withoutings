package ratelimit

import (
	"context"
	"time"

	"golang.org/x/time/rate"
)

var DefaultWithingsAPIRateLimiter = NewLeakyBucket(120 * rate.Every(time.Minute))

// LeakyBucket allows actions at a constant rate.
type LeakyBucket struct {
	l chan struct{}
}

func NewLeakyBucket(r rate.Limit) *LeakyBucket {
	sw := &LeakyBucket{
		l: make(chan struct{}, 1),
	}

	go sw.fillAtRate(r)

	return sw
}

func (sw *LeakyBucket) Wait(ctx context.Context) error {
	if sw == nil {
		panic("nil LeakyBucket")
	}
	select {
	case <-sw.l:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (sw *LeakyBucket) Allow() bool {
	select {
	case <-sw.l:
		go func() {
			// We "spent" a token without using it, so we need to put it back
			sw.l <- struct{}{}
		}()
		return true
	default:
		return false
	}
}

func (sw *LeakyBucket) fillAtRate(r rate.Limit) {
	period := time.Second / time.Duration(r)
	for {
		sw.l <- struct{}{}
		time.Sleep(period)
	}
}
