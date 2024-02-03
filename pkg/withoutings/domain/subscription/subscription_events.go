package subscription

import "github.com/google/uuid"

// TODO: Separate domain events from watermill events, JSON encoding, etc.

type RawNotificationReceived struct {
	RawNotificationUUID uuid.UUID
}

type NotificationReceived struct {
	NotificationUUID uuid.UUID
}

type NotificationDataFetched struct {
	NotificationUUID uuid.UUID
}

type NotificationDataFetchFailed struct {
	NotificationUUID uuid.UUID
}
