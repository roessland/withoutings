package subscription

import "github.com/google/uuid"

type RawNotificationReceived struct {
	RawNotificationUUID uuid.UUID
}

type NotificationReceived struct {
	NotificationUUID uuid.UUID
}
