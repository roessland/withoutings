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

	// Make sure access token is fresh

	// TODO add IP filtering for webhook requests
	callbackURL := h.cfg.WebsiteURL + "withings/webhooks/" + h.cfg.WithingsWebhookSecret

	params := withings.NewNotifySubscribeParams()
	params.Appli = 1
	params.Callbackurl = callbackURL
	params.Comment = "test"
	_, err = h.withingsRepo.NotifySubscribe(ctx, acc.WithingsAccessToken, params)
	if err != nil {
		return err
	}
	err = h.subscriptionRepo.CreateSubscription(ctx, subscription.NewSubscription(
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
