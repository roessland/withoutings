package command

import (
	"context"
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
	return h.subscriptionRepo.CreateSubscription(ctx, subscription.NewSubscription(
		acc.AccountID,
		1,
		"https://withings.roessland.com/withings/webhooks/"+webhookSecret,
		webhookSecret,
	))
}

func NewSubscribeAccountHandler(accountRepo account.Repo, subscriptionsRepo subscription.Repo) SubscribeAccountHandler {
	return subscribeAccountHandler{
		accountRepo:      accountRepo,
		subscriptionRepo: subscriptionsRepo,
	}
}

type subscribeAccountHandler struct {
	accountRepo      account.Repo
	subscriptionRepo subscription.Repo
}
