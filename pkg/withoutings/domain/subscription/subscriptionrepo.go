package subscription

import (
	"context"
	"errors"
	"github.com/google/uuid"
)

var ErrSubscriptionNotFound = errors.New("subscription not found")
var ErrRawNotificationNotFound = errors.New("raw notification not found")
var ErrNotificationNotFound = errors.New("notification not found")
var ErrSubscriptionAlreadyExists = errors.New("subscription for given account and appli already exists")

//go:generate mockery --name Repo --filename subscriptionrepo_mock.go
type Repo interface {
	AllNotificationCategories(ctx context.Context) ([]NotificationCategory, error)
	GetSubscriptionByUUID(ctx context.Context, subscriptionUUID uuid.UUID) (*Subscription, error)
	GetSubscriptionsByAccountUUID(ctx context.Context, accountUUID uuid.UUID) ([]*Subscription, error)
	GetSubscriptionByWebhookSecret(ctx context.Context, webhookSecret string) (*Subscription, error)
	GetPendingSubscriptions(ctx context.Context) ([]*Subscription, error)
	CreateSubscriptionIfNotExists(ctx context.Context, subscription *Subscription) error
	ListSubscriptions(ctx context.Context) ([]*Subscription, error)
	CreateRawNotification(ctx context.Context, rawNotification *RawNotification) error
	GetRawNotificationByUUID(ctx context.Context, rawNotificationUUID uuid.UUID) (*RawNotification, error)
	GetPendingRawNotifications(ctx context.Context) ([]*RawNotification, error)
	DeleteRawNotification(ctx context.Context, rawNotification *RawNotification) error
	Update(ctx context.Context, subscriptionUUID uuid.UUID, updateFn func(ctx context.Context, sub *Subscription) (*Subscription, error)) error
	UpdateRawNotification(ctx context.Context, rawNotificationUUID uuid.UUID, updateFn func(ctx context.Context, rawNotification *RawNotification) (*RawNotification, error)) error
	CreateNotification(ctx context.Context, notification *Notification) error
	GetNotificationsByAccountUUID(ctx context.Context, accountUUID uuid.UUID) ([]*Notification, error)
	GetNotificationByUUID(ctx context.Context, notificationUUID uuid.UUID) (*Notification, error)
	UpdateNotification(ctx context.Context, notificationUUID uuid.UUID, updateFn func(ctx context.Context, notification *Notification) (*Notification, error)) error
	StoreNotificationData(ctx context.Context, notificationData *NotificationData) error
	GetNotificationDataByNotificationUUID(ctx context.Context, notificationUUID uuid.UUID) ([]*NotificationData, error)
}
