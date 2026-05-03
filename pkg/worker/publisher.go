package worker

import "github.com/ThreeDotsLabs/watermill/message"

// noopClosePublisher delegates Publish to its inner publisher and makes Close
// a no-op. Used to hand the same shared publisher to multiple watermill
// handlers without each one independently closing it on shutdown.
//
// Watermill's router runs `defer h.publisher.Close()` per handler when the
// handler exits (router.go ~625 in v1.3.x – v1.5.x). With one shared
// publisher and N handlers, that's N concurrent Close calls on the same
// underlying publisher, which races on watermill-sql's unsynchronised
// `closed` flag and risks a double close-of-channel panic. Wrapping the
// publisher per handler with a no-op Close means the worker's deferred
// single Close is the only one that actually runs.
type noopClosePublisher struct {
	inner message.Publisher
}

func (p noopClosePublisher) Publish(topic string, msgs ...*message.Message) error {
	return p.inner.Publish(topic, msgs...)
}

func (noopClosePublisher) Close() error { return nil }
