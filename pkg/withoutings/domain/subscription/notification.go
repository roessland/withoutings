package subscription

import (
	"errors"
	"github.com/google/uuid"
	"net/url"
	"time"
)

type NotificationParams struct {
	UserID    string
	StartDate string
	EndDate   string
	Appli     string
}

func NotificationParamsFromString(s string) (NotificationParams, error) {
	var p NotificationParams
	params, err := url.ParseQuery(s)
	if err != nil {
		return p, err
	}
	for paramKey, paramVal := range params {
		if len(paramVal) != 1 {
			continue
		}
		if paramKey == "userid" {
			p.UserID = paramVal[0]
		} else if paramKey == "startdate" {
			p.StartDate = paramVal[0]
		} else if paramKey == "enddate" {
			p.EndDate = paramVal[0]
		} else if paramKey == "appli" {
			p.Appli = paramVal[0]
		}
	}
	return p, err
}

// Notification represents a notification received from the Withings API,
// after it has been ingested and processed. It corresponds to a single user.
// Any raw notifications that do not correspond to a user will never become a Notification.
type Notification struct {
	notificationUUID uuid.UUID

	// accountUUID is the UUID of the account that received the notification.
	accountUUID uuid.UUID

	// receivedAt is the time the notification was received.
	receivedAt time.Time

	// params are the query parameters of the POST request from Withings.
	// Example: {UserID: "133337", StartDate: "1682809920", EndDate: "1682837220", Appli: "44"}
	params string

	// data is the response body from the Withings API using the provided query parameters.
	data []byte

	// fetchedAt is the time the data was fetched from the Withings API.
	fetchedAt time.Time

	// rawNotificationUUID is the UUID of the RawNotification that was processed to create this Notification.
	rawNotificationUUID uuid.UUID

	// source is the IP address of the Withings API server that sent the notification.
	source string
}

type NewNotificationParams struct {
	NotificationUUID    uuid.UUID
	AccountUUID         uuid.UUID
	ReceivedAt          time.Time
	Params              string
	Data                []byte
	FetchedAt           time.Time
	RawNotificationUUID uuid.UUID
	Source              string
}

// NewNotification validates input and returns a new RawNotification
func NewNotification(p NewNotificationParams) (Notification, error) {
	if p.NotificationUUID == uuid.Nil {
		return Notification{}, errors.New("notificationUUID cannot be nil")
	}

	if p.AccountUUID == uuid.Nil {
		return Notification{}, errors.New("accountUUID cannot be nil")
	}

	if p.ReceivedAt.IsZero() {
		return Notification{}, errors.New("zero receivedAt")
	}

	if p.FetchedAt.IsZero() {
		return Notification{}, errors.New("zero fetchedAt")
	}

	return Notification{
		notificationUUID:    p.NotificationUUID,
		accountUUID:         p.AccountUUID,
		receivedAt:          p.ReceivedAt,
		params:              p.Params,
		data:                p.Data,
		fetchedAt:           p.FetchedAt,
		rawNotificationUUID: p.RawNotificationUUID,
		source:              p.Source,
	}, nil
}

// UUID returns the UUID.
func (r Notification) UUID() uuid.UUID {
	return r.rawNotificationUUID
}

// AccountUUID returns the UUID of the account that received the notification.
func (r Notification) AccountUUID() uuid.UUID {
	return r.accountUUID
}

// ReceivedAt returns the time the notification was received.
func (r Notification) ReceivedAt() time.Time {
	return r.receivedAt
}

// Params returns the query parameters of the POST request from Withings.
// Example:  "userid=133337&startdate=1682809920&enddate=1682837220&appli=44"
func (r Notification) Params() string {
	return r.params
}

// Data returns the response body from the Withings API using the provided query parameters.
func (r Notification) Data() []byte {
	return r.data
}

// FetchedAt returns the time the data was fetched from the Withings API.
func (r Notification) FetchedAt() time.Time {
	return r.fetchedAt
}

// RawNotificationUUID returns the UUID of the RawNotification that was processed to create this Notification.
func (r Notification) RawNotificationUUID() uuid.UUID {
	return r.rawNotificationUUID
}

// Source returns the source.
func (r Notification) Source() string {
	return r.source
}
