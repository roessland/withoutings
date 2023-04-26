package subscription_test

import (
	"github.com/google/uuid"
	"github.com/roessland/withoutings/pkg/db"
	"github.com/roessland/withoutings/pkg/testctx"
	"github.com/roessland/withoutings/pkg/testdb"
	subscriptionAdapter "github.com/roessland/withoutings/pkg/withoutings/adapter/subscription"
	"github.com/roessland/withoutings/pkg/withoutings/domain/subscription"
	"github.com/stretchr/testify/require"
	"testing"
)

var _ subscription.Repo = subscriptionAdapter.PgRepo{}

func TestSubscriptionPgRepo_CreateSubscriptionIfNotExists(t *testing.T) {
	ctx := testctx.New()
	database := testdb.New(ctx)
	defer database.Drop(ctx)
	queries := db.New(database)
	repo := subscriptionAdapter.NewPgRepo(database.Pool, queries)

	// Subscription has mandatory foreign key to account.
	withingsUserID := uuid.NewString()
	accountUUID := uuid.New()
	err := queries.CreateAccount(ctx, db.CreateAccountParams{
		AccountUuid:    accountUUID,
		WithingsUserID: withingsUserID,
	})
	require.NoError(t, err)

	t.Run("CreateSubscriptionIfNotExists creates subscription", func(t *testing.T) {
		// domain object
		sub := subscription.NewSubscription(
			uuid.New(),
			accountUUID,
			2,
			"https://yolo.com/",
			"comment",
			"webhooksecret",
			subscription.StatusActive,
		)

		// create in DB
		err = repo.CreateSubscriptionIfNotExists(ctx, sub)
		require.NoError(t, err)

		// retrieve inserted object
		insertedSub, err := queries.GetSubscriptionByAccountUUIDAndAppli(ctx,
			db.GetSubscriptionByAccountUUIDAndAppliParams{
				AccountUuid: accountUUID,
				Appli:       2,
			})
		require.NoError(t, err)
		require.EqualValues(t, "https://yolo.com/", insertedSub.Callbackurl)
		require.EqualValues(t, "comment", insertedSub.Comment)
		require.EqualValues(t, subscription.StatusActive, insertedSub.Status)

		// insert same object again, should be error
		err = repo.CreateSubscriptionIfNotExists(ctx, sub)
		require.Error(t, err)
		require.ErrorIs(t, err, subscription.ErrSubscriptionAlreadyExists)
	})
}
