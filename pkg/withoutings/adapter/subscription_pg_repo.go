package adapter

import (
	"context"
	"github.com/jackc/pgx/v4"
	"github.com/pkg/errors"
	"github.com/roessland/withoutings/pkg/repos/db"
	"github.com/roessland/withoutings/pkg/withoutings/domain/subscription"
)

type SubscriptionPgRepo struct {
	queries *db.Queries
}

func NewSubscriptionPgRepo(queries *db.Queries) SubscriptionPgRepo {
	return SubscriptionPgRepo{
		queries: queries,
	}
}

func (r SubscriptionPgRepo) GetSubscriptionByID(ctx context.Context, subscriptionID int64) (subscription.Subscription, error) {
	dbSub, err := r.queries.GetSubscriptionByID(ctx, subscriptionID)
	if err == pgx.ErrNoRows {
		return subscription.Subscription{}, subscription.NotFoundError{}
	}
	if err != nil {
		return subscription.Subscription{}, errors.Wrap(err, "unable to retrieve subscription")
	}
	return toDomainSubscription(dbSub), nil
}

func (r SubscriptionPgRepo) GetSubscriptionsByAccountID(ctx context.Context, accountID int64) ([]subscription.Subscription, error) {
	var subscriptions []subscription.Subscription
	dbSubscriptions, err := r.queries.GetSubscriptionsByAccountID(ctx, accountID)
	if err != nil {
		return nil, err
	}
	for _, dbSub := range dbSubscriptions {
		subscriptions = append(subscriptions, toDomainSubscription(dbSub))
	}
	return subscriptions, nil
}

func (r SubscriptionPgRepo) GetSubscriptionByWebhookSecret(ctx context.Context, webhookSecret string) (subscription.Subscription, error) {
	dbSub, err := r.queries.GetSubscriptionByWebhookSecret(ctx, webhookSecret)
	if err == pgx.ErrNoRows {
		return subscription.Subscription{}, subscription.NotFoundError{}
	}
	if err != nil {
		return subscription.Subscription{}, errors.Wrap(err, "unable to retrieve subscription")
	}
	return toDomainSubscription(dbSub), nil
}

func (r SubscriptionPgRepo) CreateSubscription(ctx context.Context, sub subscription.Subscription) error {
	return r.queries.CreateSubscription(ctx, db.CreateSubscriptionParams{
		AccountID:     sub.AccountID,
		Appli:         int32(sub.Appli),
		Callbackurl:   sub.CallbackURL,
		WebhookSecret: sub.WebhookSecret,
		Comment:       sub.Comment,
		Status:        string(sub.Status),
	})
}

func (r SubscriptionPgRepo) ListSubscriptions(ctx context.Context) ([]subscription.Subscription, error) {
	var subscriptions []subscription.Subscription
	dbSubscriptions, err := r.queries.ListSubscriptions(ctx)
	if err != nil {
		return nil, err
	}
	for _, dbSub := range dbSubscriptions {
		subscriptions = append(subscriptions, toDomainSubscription(dbSub))
	}
	return subscriptions, nil
}

func toDomainSubscription(dbSub db.Subscription) subscription.Subscription {
	return subscription.Subscription{
		SubscriptionID: dbSub.SubscriptionID,
		AccountID:      dbSub.AccountID,
		Appli:          int(dbSub.Appli),
		CallbackURL:    dbSub.Callbackurl,
		Comment:        dbSub.Comment,
	}
}

func (r SubscriptionPgRepo) CreateRawNotification(ctx context.Context, rawNotification subscription.RawNotification) error {
	return r.queries.CreateRawNotification(ctx, db.CreateRawNotificationParams{
		Source: rawNotification.Source,
		Status: string(rawNotification.Status),
		Data:   rawNotification.Data,
	})
}

func (r SubscriptionPgRepo) GetPendingRawNotifications(ctx context.Context) ([]subscription.RawNotification, error) {
	var rawNotifications []subscription.RawNotification
	dbRawNotifications, err := r.queries.GetPendingRawNotifications(ctx)
	if err != nil {
		return nil, err
	}
	for _, dbRawNotification := range dbRawNotifications {
		rawNotifications = append(rawNotifications, subscription.RawNotification{
			RawNotificationID: dbRawNotification.RawNotificationID,
			Source:            dbRawNotification.Source,
			Data:              dbRawNotification.Data,
			Status:            subscription.RawNotificationStatus(dbRawNotification.Status),
		})
	}
	return rawNotifications, nil
}
