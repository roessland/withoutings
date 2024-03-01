package subscription_test

import (
	"github.com/google/uuid"
	"github.com/roessland/withoutings/pkg/withoutings/domain/subscription"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestMustNewNotification(t *testing.T) {
	nUUID := uuid.New()
	n := subscription.MustNewNotification(subscription.NewNotificationParams{
		NotificationUUID:    nUUID,
		AccountUUID:         uuid.New(),
		ReceivedAt:          time.Now(),
		Params:              "",
		DataStatus:          subscription.NotificationDataStatusAwaitingFetch,
		Data:                nil,
		FetchedAt:           nil,
		RawNotificationUUID: uuid.New(),
		Source:              "",
	})

	assert.Equal(t, nUUID, n.UUID(), "UUID getter should return the UUID")
}
