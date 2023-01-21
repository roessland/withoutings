package adapter_test

import (
	"github.com/roessland/withoutings/pkg/withoutings/adapter"
	"github.com/roessland/withoutings/pkg/withoutings/domain/subscription"
)

var _ subscription.Repo = adapter.SubscriptionPgRepo{}
