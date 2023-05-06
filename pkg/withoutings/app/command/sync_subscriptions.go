package command

import (
	"context"
	"github.com/google/uuid"
	"github.com/roessland/withoutings/pkg/config"
	"github.com/roessland/withoutings/pkg/logging"
	"github.com/roessland/withoutings/pkg/web/templates"
	withingsSvc "github.com/roessland/withoutings/pkg/withoutings/app/service/withings"
	"github.com/roessland/withoutings/pkg/withoutings/domain/account"
	"github.com/roessland/withoutings/pkg/withoutings/domain/subscription"
	"github.com/roessland/withoutings/pkg/withoutings/domain/withings"
)

type SyncSubscriptions struct {
	Account *account.Account
}

type SyncSubscriptionsHandler interface {
	Handle(ctx context.Context, cmd SyncSubscriptions) error
}

func (h syncSubscriptionsHandler) Handle(ctx context.Context, cmd SyncSubscriptions) error {
	log := logging.MustGetLoggerFromContext(ctx)

	// Check WithingsRepo for each notification category.
	categories, err := h.subscriptionRepo.AllNotificationCategories(ctx)
	if err != nil {
		log.WithError(err).WithField("event", "error.subscriptionRepo.AllNotificationCategories.failed").Error()
		return err
	}

	for _, cat := range categories {
		notifyListResponse, err := h.withingsSvc.NotifyList(ctx, cmd.Account, withings.NewNotifyListParams(cat.Appli))
		if err != nil {
			log.WithError(err).WithField("event", "error.Withings.NotifyList.failed").Error()
			return err
		}
		if len(notifyListResponse.Body.Profiles) == 0 {
			withingsSubscriptions = append(withingsSubscriptions, templates.SubscriptionsWithingsPageItem{
				Appli:            cat.Appli,
				AppliDescription: cat.Description,
				Exists:           false,
			})
		}
		for _, profile := range notifyListResponse.Body.Profiles {
			withingsSubscriptions = append(withingsSubscriptions, templates.SubscriptionsWithingsPageItem{
				Appli:            profile.Appli,
				AppliDescription: cat.Description,
				Exists:           true,
				Comment:          profile.Comment,
			})
		}
	}

	// TODO: Make sure access token is fresh

	// TODO: add IP filtering for webhook requests

	// Subscribe
	params := withings.NewNotifySubscribeParams()
	params.Appli = cmd.Appli
	callbackURL := h.cfg.WebsiteURL + "withings/webhooks/" + h.cfg.WithingsWebhookSecret
	params.Callbackurl = callbackURL
	params.Comment = "test"
	_, err = h.withingsRepo.NotifySubscribe(ctx, acc.WithingsAccessToken(), params)
	if err != nil {
		return err
	}

	// Save subscription
	err = h.subscriptionRepo.CreateSubscriptionIfNotExists(ctx, subscription.NewSubscription(
		uuid.New(),
		acc.UUID(),
		params.Appli,
		callbackURL,
		"test",
		h.cfg.WithingsWebhookSecret,
		subscription.StatusActive,
	))
	if err != nil {
		return err
	}

	return nil
}

func NewSyncSubscriptionsHandler(
	accountRepo account.Repo,
	subscriptionsRepo subscription.Repo,
	withingsSvc withingsSvc.Service,
	cfg *config.Config,
) SyncSubscriptionsHandler {
	return syncSubscriptionsHandler{
		subscriptionRepo: subscriptionsRepo,
		withingsSvc:      withingsSvc,
		cfg:              cfg,
	}
}

type syncSubscriptionsHandler struct {
	subscriptionRepo subscription.Repo
	withingsSvc      withingsSvc.Service
	cfg              *config.Config
}
