package command

import (
	"context"
	"github.com/roessland/withoutings/pkg/withoutings/clients/withingsapi"
	"github.com/roessland/withoutings/pkg/withoutings/domain/account"
	"github.com/roessland/withoutings/pkg/withoutings/domain/subscription"
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

	webhookSecret := subscription.RandomWebhookSecret()
	callbackURL := "https://withings.roessland.com/withings/webhooks/" + webhookSecret

	params := withingsapi.NewNotifySubscribeParams()
	params.Appli = 1
	params.Callbackurl = callbackURL
	params.Comment = "test"
	_, err = h.withingsClient.WithAccessToken(acc.WithingsAccessToken).NotifySubscribe(ctx, params)
	if err != nil {
		return err
	}
	err = h.subscriptionRepo.CreateSubscription(ctx, subscription.NewSubscription(
		acc.AccountID,
		params.Appli,
		callbackURL,
		"test",
		webhookSecret,
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
	withingsClient *withingsapi.Client,
) SubscribeAccountHandler {
	return subscribeAccountHandler{
		accountRepo:      accountRepo,
		subscriptionRepo: subscriptionsRepo,
		withingsClient:   withingsClient,
	}
}

type subscribeAccountHandler struct {
	accountRepo      account.Repo
	subscriptionRepo subscription.Repo
	withingsClient   *withingsapi.Client
}
