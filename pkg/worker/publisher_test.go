package worker

import (
	"sync/atomic"
	"testing"

	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/stretchr/testify/assert"
)

// recordingPublisher counts Publish/Close so the wrapper test can verify
// Publish delegates and Close is a no-op.
type recordingPublisher struct {
	publishes atomic.Int32
	closes    atomic.Int32
}

func (p *recordingPublisher) Publish(string, ...*message.Message) error {
	p.publishes.Add(1)
	return nil
}

func (p *recordingPublisher) Close() error {
	p.closes.Add(1)
	return nil
}

func TestNoopClosePublisher_DelegatesPublishButNotClose(t *testing.T) {
	inner := &recordingPublisher{}
	wrapped := noopClosePublisher{inner: inner}

	assert.NoError(t, wrapped.Publish("topic", message.NewMessage("1", nil)))
	assert.Equal(t, int32(1), inner.publishes.Load(), "Publish must delegate")

	// Watermill's per-handler shutdown calls Close on each handler's
	// publisher. With many handlers sharing the inner publisher, the worker
	// hands each a fresh wrapper; their Close calls must not propagate.
	for i := 0; i < 10; i++ {
		assert.NoError(t, wrapped.Close())
	}
	assert.Equal(t, int32(0), inner.closes.Load(), "Close must be a no-op")
}
