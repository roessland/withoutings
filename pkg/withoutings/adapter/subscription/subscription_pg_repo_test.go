package subscription_test

import (
	subscription2 "github.com/roessland/withoutings/pkg/withoutings/adapter/subscription"
	"github.com/roessland/withoutings/pkg/withoutings/domain/subscription"
)

var _ subscription.Repo = subscription2.SubscriptionPgRepo{}
