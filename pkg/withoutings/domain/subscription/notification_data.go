package subscription

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	"time"
)

// CODE SMELL:
// Should be merged with the Notification aggregate, and persisted together.

// NotificationData represents the data that was fetched from the Withings API
// in response to a notification.
type NotificationData struct {
	notificationDataUUID uuid.UUID

	// notificationUUID is the UUID of the notification this data corresponds to.
	notificationUUID uuid.UUID

	// accountUUID is the UUID of the account that received the notification.
	accountUUID uuid.UUID

	// fetchedAt is the time the data was fetched.
	fetchedAt time.Time

	// service is the service that was called to fetch the data. E.g. "Measure - Getmeas"
	service NotificationDataService

	// data is the data fetched from the Withings API. Must be valid JSON.
	data []byte
}

type NewNotificationDataParams struct {
	NotificationDataUUID uuid.UUID
	NotificationUUID     uuid.UUID
	AccountUUID          uuid.UUID
	Service              NotificationDataService
	Data                 []byte
	FetchedAt            time.Time
}

// NewNotificationData validates input and returns a new NotificationData
func NewNotificationData(p NewNotificationDataParams) (*NotificationData, error) {
	if p.NotificationDataUUID == uuid.Nil {
		return nil, errors.New("notificationDataUUID cannot be nil")
	}

	if p.NotificationUUID == uuid.Nil {
		return nil, errors.New("notificationUUID cannot be nil")
	}

	if p.AccountUUID == uuid.Nil {
		return nil, errors.New("accountUUID cannot be nil")
	}

	if p.FetchedAt.IsZero() {
		return nil, errors.New("zero fetchedAt")
	}

	if p.Service == "" {
		return nil, errors.New("service cannot be empty")
	}

	if len(p.Data) == 0 {
		return nil, errors.New("data cannot be empty")
	}

	return &NotificationData{
		notificationDataUUID: p.NotificationDataUUID,
		notificationUUID:     p.NotificationUUID,
		accountUUID:          p.AccountUUID,
		fetchedAt:            p.FetchedAt,
		service:              p.Service,
		data:                 p.Data,
	}, nil
}

// MustNewNotificationData returns a new NotificationData or panics if it fails. For testing.
func MustNewNotificationData(p NewNotificationDataParams) *NotificationData {
	n, err := NewNotificationData(p)
	if err != nil {
		panic(err)
	}
	return n
}

func (r *NotificationData) String() string {
	return fmt.Sprintf("NotificationDataUUID: %s, AccountUUID: %s,  NotificationUUID: %s, FetchedAt: %s, Service: %s, len(Data): %d}",
		r.notificationDataUUID, r.accountUUID, r.notificationUUID, r.fetchedAt, r.service, len(r.data))
}

// UUID returns the UUID.
func (r *NotificationData) UUID() uuid.UUID {
	return r.notificationDataUUID
}

// NotificationUUID returns the UUID of the notification this data corresponds to.
func (r *NotificationData) NotificationUUID() uuid.UUID {
	return r.notificationUUID
}

// AccountUUID returns the UUID of the account that received the notification.
func (r *NotificationData) AccountUUID() uuid.UUID {
	return r.accountUUID
}

// FetchedAt returns the time the data was fetched from the Withings API.
func (r *NotificationData) FetchedAt() time.Time {
	return r.fetchedAt
}

// Service returns the service that was called to fetch the data. E.g. "Measure - Getmeas".
func (r *NotificationData) Service() NotificationDataService {
	return r.service
}

// Data returns the data fetched.
func (r *NotificationData) Data() []byte {
	return r.data
}
