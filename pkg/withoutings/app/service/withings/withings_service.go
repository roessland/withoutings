package withings

import (
	"context"
	"github.com/roessland/withoutings/pkg/withoutings/domain/account"
	"github.com/roessland/withoutings/pkg/withoutings/domain/withings"
)

type Service struct {
	repo        withings.Repo
	accountRepo account.Repo
}

func (s *Service) NotifyList(ctx context.Context, account *account.Account, params withings.NotifyListParams) (*withings.NotifyListResponse, error) {
	// TODO: Handle expired token
	return s.repo.NotifyList(ctx, account.WithingsAccessToken(), params)
}
