package command

import (
	"context"

	"github.com/google/uuid"
	"github.com/roessland/withoutings/pkg/config"
	withingsSvc "github.com/roessland/withoutings/pkg/withoutings/app/service/withings"
	"github.com/roessland/withoutings/pkg/withoutings/domain/account"
	"github.com/roessland/withoutings/pkg/withoutings/domain/subscription"
	"github.com/roessland/withoutings/pkg/withoutings/domain/withings"
)

type SubscribeAccount struct {
	Account account.Account
	Appli   int
}

type SubscribeAccountHandler interface {
	Handle(ctx context.Context, cmd SubscribeAccount) error
}

func (h subscribeAccountHandler) Handle(ctx context.Context, cmd SubscribeAccount) (err error) {
	// Ensure account exists
	acc, err := h.accountRepo.GetAccountByWithingsUserID(ctx, cmd.Account.WithingsUserID())
	if err != nil {
		return err
	}

	// TODO: Make sure access token is fresh

	// TODO: add IP filtering for webhook requests

	// Subscribe
	params := withings.NewNotifySubscribeParams()
	params.Appli = cmd.Appli
	callbackURL := h.cfg.WebsiteURL + "withings/webhooks/" + h.cfg.WithingsWebhookSecret
	params.Callbackurl = callbackURL
	params.Comment = "test"
	_, err = h.withingsSvc.NotifySubscribe(ctx, acc, params)
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

func NewSubscribeAccountHandler(
	accountRepo account.Repo,
	subscriptionsRepo subscription.Repo,
	withingsSvc withingsSvc.Service,
	cfg *config.Config,
) SubscribeAccountHandler {
	return subscribeAccountHandler{
		accountRepo:      accountRepo,
		subscriptionRepo: subscriptionsRepo,
		withingsSvc:      withingsSvc,
		cfg:              cfg,
	}
}

type subscribeAccountHandler struct {
	accountRepo      account.Repo
	subscriptionRepo subscription.Repo
	withingsSvc      withingsSvc.Service
	cfg              *config.Config
}
