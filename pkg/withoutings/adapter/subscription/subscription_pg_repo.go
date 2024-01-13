package subscription

import (
	"context"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
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

func (r PgRepo) WithTx(tx pgx.Tx) PgRepo {
	return PgRepo{
		db:      r.db,
		queries: r.queries.WithTx(tx),
	}
}

func (r PgRepo) GetSubscriptionByUUID(ctx context.Context, subscriptionUUID uuid.UUID) (*subscription.Subscription, error) {
	dbSub, err := r.queries.GetSubscriptionByUUID(ctx, subscriptionUUID)
	if err == pgx.ErrNoRows {
		return nil, subscription.ErrSubscriptionNotFound
	}
	if err != nil {
		return nil, errors.Wrap(err, "unable to retrieve subscription")
	}
	return toDomainSubscription(dbSub), nil
}

func (r PgRepo) GetSubscriptionsByAccountUUID(ctx context.Context, accountUUID uuid.UUID) ([]*subscription.Subscription, error) {
	dbSubscriptions, err := r.queries.GetSubscriptionsByAccountUUID(ctx, accountUUID)
	if err != nil {
		return nil, err
	}
	return toDomainSubscriptions(dbSubscriptions), nil
}

func (r PgRepo) GetSubscriptionByWebhookSecret(ctx context.Context, webhookSecret string) (*subscription.Subscription, error) {
	dbSub, err := r.queries.GetSubscriptionByWebhookSecret(ctx, webhookSecret)
	if err == pgx.ErrNoRows {
		return nil, subscription.ErrSubscriptionNotFound
	}
	if err != nil {
		return nil, errors.Wrap(err, "unable to retrieve subscription")
	}
	return toDomainSubscription(dbSub), nil
}

func (r PgRepo) GetPendingSubscriptions(ctx context.Context) ([]*subscription.Subscription, error) {
	dbSubscriptions, err := r.queries.GetPendingSubscriptions(ctx)
	if err != nil {
		return nil, err
	}
	return toDomainSubscriptions(dbSubscriptions), nil
}

func (r PgRepo) CreateSubscriptionIfNotExists(ctx context.Context, sub *subscription.Subscription) error {
	return r.createSubscriptionIfNotExists(ctx, sub)
}

func (r PgRepo) createSubscriptionIfNotExists(ctx context.Context, sub *subscription.Subscription) error {
	tx, err := r.db.BeginTx(ctx, pgx.TxOptions{})
	defer tx.Rollback(ctx)

	inTransaction := r.queries.WithTx(tx)

	// Check if exists
	_, err = inTransaction.GetSubscriptionByAccountUUIDAndAppli(ctx,
		db.GetSubscriptionByAccountUUIDAndAppliParams{
			AccountUuid: sub.AccountUUID(),
			Appli:       int32(sub.Appli()),
		})
	if err == nil {
		return subscription.ErrSubscriptionAlreadyExists
	}
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return err
	}

	// Doesn't exist; create one.
	err = inTransaction.CreateSubscription(ctx, db.CreateSubscriptionParams{
		SubscriptionUuid: sub.UUID(),
		AccountUuid:      sub.AccountUUID(),
		Appli:            int32(sub.Appli()),
		Callbackurl:      sub.CallbackURL(),
		WebhookSecret:    sub.WebhookSecret(),
		Comment:          sub.Comment(),
		Status:           string(sub.Status()),
	})
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func (r PgRepo) ListSubscriptions(ctx context.Context) ([]*subscription.Subscription, error) {
	dbSubscriptions, err := r.queries.ListSubscriptions(ctx)
	if err != nil {
		return nil, err
	}
	return toDomainSubscriptions(dbSubscriptions), nil
}

func toDomainSubscriptions(dbSubs []db.Subscription) []*subscription.Subscription {
	var subscriptions []*subscription.Subscription
	for _, dbSub := range dbSubs {
		subscriptions = append(subscriptions, toDomainSubscription(dbSub))
	}
	return subscriptions
}

func toDomainSubscription(dbSub db.Subscription) *subscription.Subscription {
	return subscription.NewSubscription(
		dbSub.SubscriptionUuid,
		dbSub.AccountUuid,
		int(dbSub.Appli),
		dbSub.Callbackurl,
		dbSub.Comment,
		dbSub.WebhookSecret,
		subscription.Status(dbSub.Status),
	)
}

func (r PgRepo) CreateRawNotification(ctx context.Context, rawNotification subscription.RawNotification) error {
	return r.queries.CreateRawNotification(ctx, db.CreateRawNotificationParams{
		RawNotificationUuid: rawNotification.UUID(),
		Source:              rawNotification.Source(),
		Status:              string(rawNotification.Status()),
		Data:                rawNotification.Data(),
	})
}

func (r PgRepo) GetRawNotificationByUUID(ctx context.Context, rawNotificationUUID uuid.UUID) (subscription.RawNotification, error) {
	dbRawNotification, err := r.queries.GetRawNotification(ctx, rawNotificationUUID)
	if err == pgx.ErrNoRows {
		return subscription.RawNotification{}, subscription.ErrRawNotificationNotFound
	}
	if err != nil {
		return subscription.RawNotification{}, errors.Wrap(err, "unable to retrieve raw notification")
	}
	return toDomainRawNotification(dbRawNotification), nil
}

func (r PgRepo) GetPendingRawNotifications(ctx context.Context) ([]subscription.RawNotification, error) {
	var rawNotifications []subscription.RawNotification
	dbRawNotifications, err := r.queries.GetPendingRawNotifications(ctx)
	if err != nil {
		return nil, err
	}
	for _, dbRawNotification := range dbRawNotifications {
		rawNotifications = append(rawNotifications, subscription.NewRawNotification(
			dbRawNotification.RawNotificationUuid,
			dbRawNotification.Source,
			dbRawNotification.Data,
			subscription.RawNotificationStatus(dbRawNotification.Status),
			dbRawNotification.ReceivedAt,
			dbRawNotification.ProcessedAt,
		))
	}
	return rawNotifications, nil
}

func (r PgRepo) DeleteRawNotification(ctx context.Context, rawNotification subscription.RawNotification) error {
	return r.queries.DeleteRawNotification(ctx, rawNotification.UUID())
}

func toDomainRawNotification(dbRawNotification db.RawNotification) subscription.RawNotification {
	return subscription.NewRawNotification(
		dbRawNotification.RawNotificationUuid,
		dbRawNotification.Source,
		dbRawNotification.Data,
		subscription.MustRawNotificationStatusFromString(dbRawNotification.Status),
		dbRawNotification.ReceivedAt,
		dbRawNotification.ProcessedAt,
	)
}

func (r PgRepo) AllNotificationCategories(ctx context.Context) ([]subscription.NotificationCategory, error) {
	var notificationCategories []subscription.NotificationCategory
	dbNotificationCategories, err := r.queries.AllNotificationCategories(ctx)
	if err != nil {
		return nil, err
	}
	for _, dbNotificationCategory := range dbNotificationCategories {
		notificationCategories = append(notificationCategories, subscription.NotificationCategory{
			Appli:       int(dbNotificationCategory.Appli),
			Scope:       dbNotificationCategory.Scope,
			Description: dbNotificationCategory.Description,
		})
	}
	return notificationCategories, nil
}

// Update updates a subscription in the database.
// updateFn is a function that takes the current account and returns the updated account.
// updateFn is called within a transaction, so it should not start its own transaction.
// TODO test that it returns the updated sub
func (r PgRepo) Update(ctx context.Context, subscriptionUUID uuid.UUID, updateFn func(ctx context.Context, sub *subscription.Subscription) (*subscription.Subscription, error)) error {
	tx, err := r.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	inTransaction := r.WithTx(tx)

	sub, err := inTransaction.GetSubscriptionByUUID(ctx, subscriptionUUID)
	if err != nil {
		return err
	}
	updatedSub, err := updateFn(ctx, sub)
	err = inTransaction.queries.UpdateSubscription(ctx, db.UpdateSubscriptionParams{
		SubscriptionUuid:    updatedSub.UUID(),
		AccountUuid:         updatedSub.AccountUUID(),
		Appli:               int32(updatedSub.Appli()),
		Callbackurl:         updatedSub.CallbackURL(),
		WebhookSecret:       updatedSub.WebhookSecret(),
		Comment:             updatedSub.Comment(),
		Status:              string(updatedSub.Status()),
		StatusLastCheckedAt: updatedSub.StatusLastCheckedAt(),
	})
	if err != nil {
		return err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return err
	}

	return nil
}

// UpdateRawNotification updates a raw notification in the database.
// updateFn is a function that takes the current raw notification and returns the updated raw notification.
// updateFn is called within a transaction, so it should not start its own transaction.
func (r PgRepo) UpdateRawNotification(ctx context.Context, rawNotificationUUID uuid.UUID, updateFn func(ctx context.Context, rawNotification *subscription.RawNotification) (*subscription.RawNotification, error)) error {
	tx, err := r.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	inTransaction := r.WithTx(tx)

	rawNotification, err := inTransaction.GetRawNotificationByUUID(ctx, rawNotificationUUID)
	if err != nil {
		return err
	}
	updatedRawNotification, err := updateFn(ctx, &rawNotification)
	err = inTransaction.queries.UpdateRawNotification(ctx, db.UpdateRawNotificationParams{
		RawNotificationUuid: updatedRawNotification.UUID(),
		Source:              updatedRawNotification.Source(),
		Status:              string(updatedRawNotification.Status()),
		Data:                updatedRawNotification.Data(),
		ReceivedAt:          updatedRawNotification.ReceivedAt(),
		ProcessedAt:         updatedRawNotification.ProcessedAt(),
	})
	if err != nil {
		return err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return err
	}

	return nil
}

type DbNotificationParams struct {
	UserID    string `json:"userid"`
	StartDate string `json:"startdate"`
	EndDate   string `json:"enddate"`
	Appli     string `json:"appli"`
}

func encodeDBNotificationParams(p subscription.NotificationParams) []byte {
	buf, err := json.Marshal(struct {
		UserID    string `json:"userid"`
		StartDate string `json:"startdate"`
		EndDate   string `json:"enddate"`
		Appli     string `json:"appli"`
	}{
		UserID:    p.UserID,
		StartDate: p.StartDate,
		EndDate:   p.EndDate,
		Appli:     p.Appli,
	})

	if err != nil {
		panic(err)
	}
	return buf
}

func (r PgRepo) CreateNotification(ctx context.Context, notification subscription.Notification) error {
	tx, err := r.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	inTransaction := r.WithTx(tx)

	err = inTransaction.queries.CreateNotification(ctx, db.CreateNotificationParams{
		NotificationUuid:    notification.UUID(),
		AccountUuid:         notification.AccountUUID(),
		ReceivedAt:          notification.ReceivedAt(),
		Params:              encodeDBNotificationParams(notification.Params()),
		Data:                notification.Data(),
		FetchedAt:           notification.FetchedAt(),
		RawNotificationUuid: notification.RawNotificationUUID(),
		Source:              notification.Source(),
	})
	if err != nil {
		return err
	}

	err = inTransaction.UpdateRawNotification(ctx, notification.RawNotificationUUID(), func(ctx context.Context, rawNotification *subscription.RawNotification) (*subscription.RawNotification, error) {
		rawNotification.MarkAsProcessed()
		return rawNotification, nil
	})

	err = tx.Commit(ctx)
	if err != nil {
		return err
	}
	return nil
}
