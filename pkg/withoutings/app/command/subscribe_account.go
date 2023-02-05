package command

import (
	"context"
	"github.com/roessland/withoutings/pkg/config"
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
	acc, err := h.accountRepo.GetAccountByWithingsUserID(ctx, cmd.Account.WithingsUserID)
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
	_, err = h.withingsRepo.NotifySubscribe(ctx, acc.WithingsAccessToken, params)
	if err != nil {
		return err
	}

	// Save subscription
	err = h.subscriptionRepo.CreateSubscriptionIfNotExists(ctx, subscription.NewSubscription(
		acc.AccountID,
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
	withingsRepo withings.Repo,
	cfg *config.Config,
) SubscribeAccountHandler {
	return subscribeAccountHandler{
		accountRepo:      accountRepo,
		subscriptionRepo: subscriptionsRepo,
		withingsRepo:     withingsRepo,
		cfg:              cfg,
	}
}

type subscribeAccountHandler struct {
	accountRepo      account.Repo
	subscriptionRepo subscription.Repo
	withingsRepo     withings.Repo
	cfg              *config.Config
}
