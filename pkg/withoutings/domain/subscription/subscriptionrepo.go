package subscription

import (
	"context"
	"fmt"
)

type NotFoundError struct {
	SubscriptionID int64
}

func (e NotFoundError) Error() string {
	return fmt.Sprintf("subscription with ID %d not found", e.SubscriptionID)
}

type Repo interface {
	GetSubscriptionByID(ctx context.Context, subscriptionID int64) (Subscription, error)
	GetSubscriptionsByAccountID(ctx context.Context, accountID int64) ([]Subscription, error)
	CreateSubscription(ctx context.Context, subscription Subscription) error
	ListSubscriptions(ctx context.Context) ([]Subscription, error)
}
