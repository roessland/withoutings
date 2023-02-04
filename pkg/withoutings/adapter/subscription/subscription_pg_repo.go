package subscription

import (
	"context"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/pkg/errors"
	"github.com/roessland/withoutings/pkg/db"
	"github.com/roessland/withoutings/pkg/withoutings/domain/subscription"
)

type PgRepo struct {
	db      *pgxpool.Pool
	queries *db.Queries
}

func NewPgRepo(db *pgxpool.Pool, queries *db.Queries) PgRepo {
	return PgRepo{
		db:      db,
		queries: queries,
	}
}

func (r PgRepo) GetSubscriptionByID(ctx context.Context, subscriptionID int64) (subscription.Subscription, error) {
	dbSub, err := r.queries.GetSubscriptionByID(ctx, subscriptionID)
	if err == pgx.ErrNoRows {
		return subscription.Subscription{}, subscription.NotFoundError{}
	}
	if err != nil {
		return subscription.Subscription{}, errors.Wrap(err, "unable to retrieve subscription")
	}
	return toDomainSubscription(dbSub), nil
}

func (r PgRepo) GetSubscriptionsByAccountID(ctx context.Context, accountID int64) ([]subscription.Subscription, error) {
	dbSubscriptions, err := r.queries.GetSubscriptionsByAccountID(ctx, accountID)
	if err != nil {
		return nil, err
	}
	return toDomainSubscriptions(dbSubscriptions), nil
}

func (r PgRepo) GetSubscriptionByWebhookSecret(ctx context.Context, webhookSecret string) (subscription.Subscription, error) {
	dbSub, err := r.queries.GetSubscriptionByWebhookSecret(ctx, webhookSecret)
	if err == pgx.ErrNoRows {
		return subscription.Subscription{}, subscription.NotFoundError{}
	}
	if err != nil {
		return subscription.Subscription{}, errors.Wrap(err, "unable to retrieve subscription")
	}
	return toDomainSubscription(dbSub), nil
}

func (r PgRepo) GetPendingSubscriptions(ctx context.Context) ([]subscription.Subscription, error) {
	dbSubscriptions, err := r.queries.GetPendingSubscriptions(ctx)
	if err != nil {
		return nil, err
	}
	return toDomainSubscriptions(dbSubscriptions), nil
}

func (r PgRepo) CreateSubscription(ctx context.Context, sub subscription.Subscription) error {
	return r.queries.CreateSubscription(ctx, db.CreateSubscriptionParams{
		AccountID:     sub.AccountID,
		Appli:         int32(sub.Appli),
		Callbackurl:   sub.CallbackURL,
		WebhookSecret: sub.WebhookSecret,
		Comment:       sub.Comment,
		Status:        string(sub.Status),
	})
}

func (r PgRepo) ListSubscriptions(ctx context.Context) ([]subscription.Subscription, error) {
	dbSubscriptions, err := r.queries.ListSubscriptions(ctx)
	if err != nil {
		return nil, err
	}
	return toDomainSubscriptions(dbSubscriptions), nil
}

func toDomainSubscriptions(dbSubs []db.Subscription) []subscription.Subscription {
	var subscriptions []subscription.Subscription
	for _, dbSub := range dbSubs {
		subscriptions = append(subscriptions, toDomainSubscription(dbSub))
	}
	return subscriptions
}

func toDomainSubscription(dbSub db.Subscription) subscription.Subscription {
	return subscription.Subscription{
		SubscriptionID: dbSub.SubscriptionID,
		AccountID:      dbSub.AccountID,
		Appli:          int(dbSub.Appli),
		CallbackURL:    dbSub.Callbackurl,
		Comment:        dbSub.Comment,
		Status:         subscription.Status(dbSub.Status),
	}
}

func (r PgRepo) CreateRawNotification(ctx context.Context, rawNotification subscription.RawNotification) error {
	return r.queries.CreateRawNotification(ctx, db.CreateRawNotificationParams{
		Source: rawNotification.Source,
		Status: string(rawNotification.Status),
		Data:   rawNotification.Data,
	})
}

func (r PgRepo) GetPendingRawNotifications(ctx context.Context) ([]subscription.RawNotification, error) {
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
