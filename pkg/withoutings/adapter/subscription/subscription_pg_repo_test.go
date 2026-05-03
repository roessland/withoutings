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
	"time"
)

var _ subscription.Repo = subscriptionAdapter.PgRepo{}

func TestSubscriptionPgRepo(t *testing.T) {
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

	t.Run("create works, get works, create duplicate fails", func(t *testing.T) {
		sub := subscription.NewSubscription(
			uuid.New(),
			accountUUID,
			2,
			"https://yolo.com/",
			"comment",
			"webhooksecret",
			subscription.StatusActive,
		)

		// Create it
		err = repo.CreateSubscriptionIfNotExists(ctx, sub)
		require.NoError(t, err)

		// Retrieve it and ensure it's the same, and has default values from Postgres.
		insertedSub, err := repo.GetSubscriptionByUUID(ctx, sub.UUID())
		require.NoError(t, err)
		require.EqualValues(t, "https://yolo.com/", insertedSub.CallbackURL())
		require.EqualValues(t, "comment", insertedSub.Comment())
		require.EqualValues(t, "webhooksecret", insertedSub.WebhookSecret())
		require.EqualValues(t, subscription.StatusActive, insertedSub.Status())
		require.True(t, insertedSub.StatusShouldBeChecked())

		// Insert the same object again, ensure it fails
		err = repo.CreateSubscriptionIfNotExists(ctx, sub)
		require.Error(t, err)
		require.ErrorIs(t, err, subscription.ErrSubscriptionAlreadyExists)
	})

	t.Run("CreateRawNotification creates and DeleteRawNotification deletes", func(t *testing.T) {
		rawNotification := subscription.NewRawNotification(
			uuid.New(),
			"ip=123.123.123.123",
			"appli=1337&foo=bar",
			subscription.RawNotificationStatusPending,
			time.Now(),
			nil,
		)

		// Create it
		err = repo.CreateRawNotification(ctx, rawNotification)
		require.NoError(t, err)

		// Retrieve it and ensure it's the same
		insertedRawNotification, err := repo.GetRawNotificationByUUID(ctx, rawNotification.UUID())
		require.NoError(t, err)
		require.EqualValues(t, rawNotification.UUID(), insertedRawNotification.UUID())
		require.EqualValues(t, rawNotification.Source(), insertedRawNotification.Source())
		require.EqualValues(t, rawNotification.Status(), insertedRawNotification.Status())
		require.EqualValues(t, rawNotification.Data(), insertedRawNotification.Data())

		// Delete it
		err = repo.DeleteRawNotification(ctx, rawNotification)

		// Ensure it was deleted
		_, err = repo.GetRawNotificationByUUID(ctx, rawNotification.UUID())
		require.Error(t, err)
	})

	t.Run("GetNotificationsByAccountUUID works", func(t *testing.T) {
		notification :=
			subscription.MustNewNotification(subscription.NewNotificationParams{
				NotificationUUID:    uuid.New(),
				AccountUUID:         accountUUID,
				ReceivedAt:          time.Now(),
				Params:              "yolo",
				DataStatus:          subscription.NotificationDataStatusAwaitingFetch,
				FetchedAt:           nil,
				RawNotificationUUID: uuid.New(),
				Source:              "",
			})
		err := repo.CreateNotification(ctx, notification)
		require.NoError(t, err)

		notifications, err := repo.GetNotificationsByAccountUUID(ctx, accountUUID)
		require.NoError(t, err)
		require.Len(t, notifications, 1)
	})

	t.Run("GetNotificationDataByAccountAndServiceAndOverlappingWindow returns only overlapping rows", func(t *testing.T) {
		// Three Sleep v2 - Get rows: only the middle one's body.series
		// overlaps the requested [3000, 4000] window. The other two are
		// strictly before / strictly after and must be filtered in SQL.
		insert := func(notifUUID uuid.UUID, body string) {
			t.Helper()
			notif := subscription.MustNewNotification(subscription.NewNotificationParams{
				NotificationUUID:    notifUUID,
				AccountUUID:         accountUUID,
				ReceivedAt:          time.Now(),
				Params:              "x",
				DataStatus:          subscription.NotificationDataStatusFetched,
				FetchedAt:           ptrTimeAdapter(time.Now()),
				RawNotificationUUID: uuid.New(),
				Source:              "test",
			})
			require.NoError(t, repo.CreateNotification(ctx, notif))
			data := subscription.MustNewNotificationData(subscription.NewNotificationDataParams{
				NotificationDataUUID: uuid.New(),
				NotificationUUID:     notifUUID,
				AccountUUID:          accountUUID,
				Service:              subscription.NotificationDataServiceSleepv2Get,
				Data:                 []byte(body),
				FetchedAt:            time.Now(),
			})
			require.NoError(t, repo.StoreNotificationData(ctx, data))
		}

		before := uuid.New()
		insert(before, `{"body":{"series":[{"startdate":1000,"enddate":2000,"state":1}]}}`)
		matching := uuid.New()
		insert(matching, `{"body":{"series":[{"startdate":3500,"enddate":3800,"state":1}]}}`)
		after := uuid.New()
		insert(after, `{"body":{"series":[{"startdate":5000,"enddate":6000,"state":1}]}}`)

		rows, err := repo.GetNotificationDataByAccountAndServiceAndOverlappingWindow(
			ctx, accountUUID, subscription.NotificationDataServiceSleepv2Get, 3000, 4000,
		)
		require.NoError(t, err)
		require.Len(t, rows, 1, "only the row whose segment overlaps [3000,4000] should be returned")
		require.Equal(t, matching, rows[0].NotificationUUID())
	})
}

func ptrTimeAdapter(t time.Time) *time.Time { return &t }
