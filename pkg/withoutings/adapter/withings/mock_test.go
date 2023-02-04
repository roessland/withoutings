package withings_test

import (
	withingsAdapter "github.com/roessland/withoutings/pkg/withoutings/adapter/withings"
	"github.com/roessland/withoutings/pkg/withoutings/domain/withings"
)

var _ withings.Repo = &withingsAdapter.MockClient{}
