package command

import (
	"context"
	"fmt"
	"github.com/roessland/withoutings/pkg/logging"
	withingsSvc "github.com/roessland/withoutings/pkg/withoutings/app/service/withings"
	"github.com/roessland/withoutings/pkg/withoutings/domain/account"
	"github.com/roessland/withoutings/pkg/withoutings/domain/subscription"
	"github.com/roessland/withoutings/pkg/withoutings/domain/withings"
)

type SyncRevokedSubscriptions struct {
	Account *account.Account
}

type SyncRevokedSubscriptionsHandler interface {
	Handle(ctx context.Context, cmd SyncRevokedSubscriptions) error
}

func (h syncRevokedSubscriptionsHandler) Handle(ctx context.Context, cmd SyncRevokedSubscriptions) error {
	log := logging.MustGetLoggerFromContext(ctx)

	// Existing subscriptions
	subs, err := h.subscriptionRepo.GetSubscriptionsByAccountUUID(ctx, cmd.Account.UUID())
	if err != nil {
		return fmt.Errorf("failed to get subscriptions: %w", err)
	}

	// Get list of active subscriptions from Withings for the categories in DB.
	for _, sub := range subs {

		notifyListResponse, err := h.withingsSvc.NotifyList(ctx, cmd.Account,
			withings.NewNotifyListParams(sub.Appli()),
		)
		if err != nil {
			log.WithError(err).WithField("event", "error.withings.NotifyList.failed").Error()
			return fmt.Errorf("NotifyList failed: %w", err)
		}

		subIsActive := false
		for _, profile := range notifyListResponse.Body.Profiles {
			if profile.CallbackURL == sub.CallbackURL() {
				subIsActive = true
				break
			}
		}

		if subIsActive {
			err := h.subscriptionRepo.Update(ctx, sub,
				func(ctx context.Context, sub *subscription.Subscription) (*subscription.Subscription, error) {
					sub.MarkAsCheckedAndActive()
					return sub, nil
				})
			if err != nil {
				return fmt.Errorf("failed to mark subscription as active: %w", err)
			}
		} else {
			err := h.subscriptionRepo.Update(ctx, sub,
				func(ctx context.Context, sub *subscription.Subscription) (*subscription.Subscription, error) {
					sub.MarkAsRevoked()
					return sub, nil
				})
			if err != nil {
				return fmt.Errorf("failed to mark subscription as revoked: %w", err)
			}
		}

	}

	return nil
}

func NewSyncRevokedSubscriptionsHandler(
	subscriptionsRepo subscription.Repo,
	withingsSvc *withingsSvc.Service,
) SyncRevokedSubscriptionsHandler {
	return syncRevokedSubscriptionsHandler{
		subscriptionRepo: subscriptionsRepo,
		withingsSvc:      withingsSvc,
	}
}

type syncRevokedSubscriptionsHandler struct {
	subscriptionRepo subscription.Repo
	withingsSvc      *withingsSvc.Service
}
