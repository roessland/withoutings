package db_test

import (
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/roessland/withoutings/pkg/db"
	"github.com/roessland/withoutings/pkg/testctx"
	"github.com/roessland/withoutings/pkg/testdb"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestSubscriptionQueries(t *testing.T) {
	ctx := testctx.New()
	database := testdb.New(ctx)
	defer database.Drop(ctx)

	queries := db.New(database)

	// Needed for constraint
	err := queries.CreateAccount(ctx, db.CreateAccountParams{})
	require.NoError(t, err)

	t.Run("CreateSubscription", func(t *testing.T) {
		createSubscriptionParams := db.CreateSubscriptionParams{
			AccountID:     1,
			Appli:         2,
			Callbackurl:   "https://mysite.com/w/asdf",
			WebhookSecret: "yolo",
			Comment:       "",
		}
		err := queries.CreateSubscription(ctx, createSubscriptionParams)
		require.NoError(t, err)

		subscriptions, err := queries.GetSubscriptionsByAccountID(ctx, 1)
		require.NoError(t, err)
		require.Len(t, subscriptions, 1)
		sub := subscriptions[0]
		require.EqualValues(t, 1, sub.AccountID)
		require.EqualValues(t, 2, sub.Appli)
		require.EqualValues(t, "https://mysite.com/w/asdf", sub.Callbackurl)
		require.EqualValues(t, "yolo", sub.WebhookSecret)
		require.EqualValues(t, "", sub.Comment)
	})
}
