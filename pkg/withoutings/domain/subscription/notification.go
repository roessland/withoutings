package subscription

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/roessland/withoutings/pkg/ptrof"
	"net/url"
	"strconv"
	"time"
)

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

	// dataStatus is the status of the corresponding data that should be/was fetched from the Withings API.
	dataStatus NotificationDataStatus

	// fetchedAt is the time the data was fetched from the Withings API.
	fetchedAt *time.Time

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
	DataStatus          NotificationDataStatus
	Data                []byte
	FetchedAt           *time.Time
	RawNotificationUUID uuid.UUID
	Source              string
}

// NewNotification validates input and returns a new RawNotification
func NewNotification(p NewNotificationParams) (*Notification, error) {
	if p.NotificationUUID == uuid.Nil {
		return nil, errors.New("notificationUUID cannot be nil")
	}

	if p.AccountUUID == uuid.Nil {
		return nil, errors.New("accountUUID cannot be nil")
	}

	if p.ReceivedAt.IsZero() {
		return nil, errors.New("zero receivedAt")
	}

	if p.FetchedAt != nil && p.FetchedAt.IsZero() {
		return nil, errors.New("zero fetchedAt")
	}

	if p.DataStatus == "" {
		return nil, errors.New("dataStatus cannot be empty")
	}

	if p.DataStatus == NotificationDataStatusAwaitingFetch && p.FetchedAt != nil {
		return nil, fmt.Errorf("fetchedAt must be nil when dataStatus is awaiting_fetch, but was %s", *p.FetchedAt)
	}

	if p.DataStatus == NotificationDataStatusFetched && p.FetchedAt == nil {
		return nil, errors.New("fetchedAt cannot be nil when dataStatus is fetched")
	}

	if p.Data != nil {
		if err := json.Unmarshal(p.Data, new(map[string]any)); err != nil {
			return nil, fmt.Errorf("notification data is not valid JSON: %w", err)
		}
	}

	if p.RawNotificationUUID == uuid.Nil {
		return nil, errors.New("rawNotificationUUID cannot be nil")
	}

	return &Notification{
		notificationUUID:    p.NotificationUUID,
		accountUUID:         p.AccountUUID,
		receivedAt:          p.ReceivedAt,
		params:              p.Params,
		data:                p.Data,
		fetchedAt:           p.FetchedAt,
		dataStatus:          p.DataStatus,
		rawNotificationUUID: p.RawNotificationUUID,
		source:              p.Source,
	}, nil
}

// MustNewNotification returns a new Notification or panics if it fails. For testing.
func MustNewNotification(p NewNotificationParams) *Notification {
	n, err := NewNotification(p)
	if err != nil {
		panic(err)
	}
	return n
}

func (r *Notification) String() string {
	return fmt.Sprintf("Notification{UUID: %s, AccountUUID: %s, ReceivedAt: %s, Params: %s, FetchedAt: %s, RawNotificationUUID: %s, Source: %s}",
		r.notificationUUID, r.accountUUID, r.receivedAt, r.params, r.fetchedAt, r.rawNotificationUUID, r.source)
}

// UUID returns the UUID.
func (r *Notification) UUID() uuid.UUID {
	return r.notificationUUID
}

// AccountUUID returns the UUID of the account that received the notification.
func (r *Notification) AccountUUID() uuid.UUID {
	return r.accountUUID
}

// ReceivedAt returns the time the notification was received.
func (r *Notification) ReceivedAt() time.Time {
	return r.receivedAt
}

// Params returns the query parameters of the POST request from Withings.
// Example:  "userid=133337&startdate=1682809920&enddate=1682837220&appli=44"
func (r *Notification) Params() string {
	return r.params
}

func (r *RawNotification) ParsedParams() (ParsedNotificationParams, error) {
	queryStr := r.data
	params, _ := url.ParseQuery(queryStr)
	startUnix, _ := strconv.ParseInt(params.Get("startdate"), 10, 64)
	endUnix, _ := strconv.ParseInt(params.Get("enddate"), 10, 64)
	appli, _ := strconv.Atoi(params.Get("appli"))
	date, _ := strconv.ParseInt(params.Get("date"), 10, 64)

	return ParsedNotificationParams{
		WithingsUserID: params.Get("userid"),
		StartDate:      time.Unix(startUnix, 0),
		StartDateStr:   params.Get("startdate"),
		EndDate:        time.Unix(endUnix, 0),
		EndDateStr:     params.Get("enddate"),
		Appli:          appli,
		AppliStr:       params.Get("appli"),
		Date:           time.Unix(date, 0),
		DateStr:        params.Get("date"),
		DeviceID:       params.Get("deviceid"),
	}, nil
}

// Data returns the response body from the Withings API using the provided query parameters.
func (r *Notification) Data() []byte {
	return r.data
}

// DataStatus returns the status of the corresponding data that should be/was fetched from the Withings API.
func (r *Notification) DataStatus() NotificationDataStatus {
	return r.dataStatus
}

// FetchedAt returns the time the data was fetched from the Withings API.
func (r *Notification) FetchedAt() *time.Time {
	return r.fetchedAt
}

// RawNotificationUUID returns the UUID of the RawNotification that was processed to create this Notification.
func (r *Notification) RawNotificationUUID() uuid.UUID {
	return r.rawNotificationUUID
}

// Source returns the source.
func (r *Notification) Source() string {
	return r.source
}

func (r *Notification) FetchedData(data []byte) {
	r.fetchedAt = ptrof.Time(time.Now())
	r.dataStatus = NotificationDataStatusFetched
	r.data = data
}

func (r *Notification) FetchFailed() {
	r.fetchedAt = ptrof.Time(time.Now())
	r.dataStatus = NotificationDataStatusFetchFailed
	r.data = nil
}
