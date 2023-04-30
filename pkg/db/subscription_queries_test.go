package db_test

import (
	"github.com/google/uuid"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/roessland/withoutings/pkg/db"
	"github.com/roessland/withoutings/pkg/testctx"
	"github.com/roessland/withoutings/pkg/testdb"
	"github.com/roessland/withoutings/pkg/withoutings/domain/subscription"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestSubscriptionQueries(t *testing.T) {
	ctx := testctx.New()
	database := testdb.New(ctx)
	defer database.Drop(ctx)

	queries := db.New(database)

	// Subscription has mandatory foreign key to account.
	accUUID := uuid.New()
	err := queries.CreateAccount(ctx, db.CreateAccountParams{
		AccountUuid:    accUUID,
		WithingsUserID: uuid.NewString(),
	})
	require.NoError(t, err)

	t.Run("CreateSubscriptions creates subscription", func(t *testing.T) {
		subUUID := uuid.New()
		createSubscriptionParams := db.CreateSubscriptionParams{
			SubscriptionUuid: subUUID,
			AccountUuid:      accUUID,
			Appli:            2,
			Callbackurl:      "https://mysite.com/w/asdf",
			WebhookSecret:    "yolo",
			Status:           string(subscription.StatusPending),
			Comment:          "",
		}
		err := queries.CreateSubscription(ctx, createSubscriptionParams)
		require.NoError(t, err)

		subscriptions, err := queries.GetSubscriptionsByAccountUUID(ctx, accUUID)
		require.NoError(t, err)
		require.Len(t, subscriptions, 1)
		sub := subscriptions[0]
		require.EqualValues(t, subUUID, sub.SubscriptionUuid)
		require.EqualValues(t, 2, sub.Appli)
		require.EqualValues(t, "https://mysite.com/w/asdf", sub.Callbackurl)
		require.EqualValues(t, "yolo", sub.WebhookSecret)
		require.EqualValues(t, subscription.StatusPending, sub.Status)
		require.EqualValues(t, "", sub.Comment)
		require.True(t, sub.StatusLastCheckedAt.Before(time.Now()), "StatusLastCheckedAt should be set to a time in the past")
	})
}
