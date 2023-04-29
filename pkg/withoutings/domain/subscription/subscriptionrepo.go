package subscription

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
)

type NotFoundError struct {
	SubscriptionID int64
}

func (e NotFoundError) Error() string {
	return fmt.Sprintf("subscription with ID %d not found", e.SubscriptionID)
}

var ErrSubscriptionAlreadyExists error = errors.New("subscription for given account and appli already exists")

//go:generate mockery --name Repo --filename subscriptionrepo_mock.go
type Repo interface {
	GetSubscriptionByUUID(ctx context.Context, subscriptionUUID uuid.UUID) (Subscription, error)
	GetSubscriptionsByAccountUUID(ctx context.Context, accountID uuid.UUID) ([]Subscription, error)
	GetSubscriptionByWebhookSecret(ctx context.Context, webhookSecret string) (Subscription, error)
	GetPendingSubscriptions(ctx context.Context) ([]Subscription, error)
	CreateSubscriptionIfNotExists(ctx context.Context, subscription Subscription) error
	ListSubscriptions(ctx context.Context) ([]Subscription, error)
	CreateRawNotification(ctx context.Context, rawNotification RawNotification) error
	GetPendingRawNotifications(ctx context.Context) ([]RawNotification, error)
	AllNotificationCategories(ctx context.Context) ([]NotificationCategory, error)
}
