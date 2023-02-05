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
	err := queries.CreateAccount(ctx, db.CreateAccountParams{
		WithingsUserID: withingsUserID,
	})
	require.NoError(t, err)
	acc, err := queries.GetAccountByWithingsUserID(ctx, withingsUserID)
	accountID := acc.AccountID

	t.Run("CreateSubscriptionIfNotExists creates subscription", func(t *testing.T) {
		err := repo.CreateSubscriptionIfNotExists(ctx, subscription.Subscription{
			AccountID:   accountID,
			Appli:       2,
			CallbackURL: "https://yolo.com/",
			Comment:     "comment",
			Status:      subscription.StatusActive,
		})
		require.NoError(t, err)

		insertedSub, err := queries.GetSubscriptionByAccountIDAndAppli(ctx,
			db.GetSubscriptionByAccountIDAndAppliParams{
				AccountID: accountID,
				Appli:     2,
			})
		require.NoError(t, err)
		require.EqualValues(t, "https://yolo.com/", insertedSub.Callbackurl)
		require.EqualValues(t, "comment", insertedSub.Comment)
		require.EqualValues(t, subscription.StatusActive, insertedSub.Status)

		err = repo.CreateSubscriptionIfNotExists(ctx, subscription.Subscription{
			AccountID:   accountID,
			Appli:       2,
			CallbackURL: "https://yolo.com/",
			Comment:     "comment",
			Status:      subscription.StatusActive,
		})
		require.Error(t, err)
		require.ErrorIs(t, err, subscription.ErrSubscriptionAlreadyExists)
	})
}
