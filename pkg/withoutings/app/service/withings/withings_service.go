package withings

import (
	"context"
	"fmt"
	"github.com/roessland/withoutings/pkg/withoutings/domain/account"
	"github.com/roessland/withoutings/pkg/withoutings/domain/withings"
)

//go:generate mockery --name Service --filename withings_service_mock.go
type Service interface {
	NotifyList(ctx context.Context, acc *account.Account, params withings.NotifyListParams) (*withings.NotifyListResponse, error)
	NotifySubscribe(ctx context.Context, acc *account.Account, params withings.NotifySubscribeParams) (*withings.NotifySubscribeResponse, error)
	SleepGetsummary(ctx context.Context, acc *account.Account, params withings.SleepGetsummaryParams) (*withings.SleepGetsummaryResponse, error)
	MeasureGetmeas(ctx context.Context, acc *account.Account, params withings.MeasureGetmeasParams) (*withings.MeasureGetmeasResponse, error)
}

type service struct {
	repo        withings.Repo
	accountRepo account.Repo
}

func NewService(withingsRepo withings.Repo, accountRepo account.Repo) Service {
	return &service{
		repo:        withingsRepo,
		accountRepo: accountRepo,
	}
}

// ensureTokenStillValid refreshes the access token for the given account and persists it.
// If someone else already updated the access token, the input account is updated with those values.
func (s *service) ensureTokenStillValid(ctx context.Context, acc *account.Account) error {
	// Token still valid, no need to refresh.
	if !acc.CanRefreshAccessToken() {
		return nil
	}

	// Get new token from Withings API
	token, err := s.repo.RefreshAccessToken(ctx, acc.WithingsRefreshToken())
	refreshSucceeded := err == nil

	if refreshSucceeded {
		// In transaction
		err = s.accountRepo.Update(ctx, acc.UUID(), func(ctx context.Context, accLatest *account.Account) (*account.Account, error) {
			// Store updated token in database
			err := accLatest.UpdateWithingsToken(token.AccessToken, token.RefreshToken, token.Expiry, token.Scope)
			if err != nil {
				return nil, fmt.Errorf("UpdateWithingsToken for latest account failed: %w", err)
			}

			// Update input account with latest values from database so we don't have to fetch the account again.
			err = acc.UpdateWithingsToken(
				accLatest.WithingsAccessToken(),
				accLatest.WithingsRefreshToken(),
				accLatest.WithingsAccessTokenExpiry(),
				accLatest.WithingsScopes(),
			)
			if err != nil {
				return nil, fmt.Errorf("UpdateWithingsToken for input account failed: %w", err)
			}
			return accLatest, nil
		})
		if err != nil {
			return fmt.Errorf("failed to update account with new token: %w", err)
		}
	}
	return nil
}

func executeWithRetry[P any, R any](s *service, fn func(ctx context.Context, accessToken string, params P) (*R, error), ctx context.Context, acc *account.Account, params P) (*R, error) {
	// To verify the logic works the when first attempt fails, this is commented out for now.
	//err := s.ensureTokenStillValid(ctx, acc)
	//if err != nil {
	//	return nil, fmt.Errorf("failed to ensure valid access token (guard): %w", err)
	//}

	// First attempt
	resp, err := fn(ctx, acc.WithingsAccessToken(), params)
	if err == nil {
		return resp, nil
	}
	if err != withings.ErrInvalidToken {
		return nil, fmt.Errorf("unexpected error from Withings API: %w", err)
	}

	// If token was invalid, refresh token
	err = s.ensureTokenStillValid(ctx, acc)
	if err != nil {
		return nil, fmt.Errorf("failed to ensure valid access token (retry): %w", err)
	}

	// Second attempt
	resp, err = fn(ctx, acc.WithingsAccessToken(), params)
	if err != nil {
		return nil, fmt.Errorf("failed to execute Withings API request after retry: %w", err)
	}
	return resp, nil
}

func (s *service) NotifyList(ctx context.Context, acc *account.Account, params withings.NotifyListParams) (*withings.NotifyListResponse, error) {
	return executeWithRetry(s, s.repo.NotifyList, ctx, acc, params)
}

func (s *service) NotifySubscribe(ctx context.Context, acc *account.Account, params withings.NotifySubscribeParams) (*withings.NotifySubscribeResponse, error) {
	return executeWithRetry(s, s.repo.NotifySubscribe, ctx, acc, params)
}

func (s *service) SleepGetsummary(ctx context.Context, acc *account.Account, params withings.SleepGetsummaryParams) (*withings.SleepGetsummaryResponse, error) {
	return executeWithRetry(s, s.repo.SleepGetsummary, ctx, acc, params)
}

func (s *service) MeasureGetmeas(ctx context.Context, acc *account.Account, params withings.MeasureGetmeasParams) (*withings.MeasureGetmeasResponse, error) {
	return executeWithRetry(s, s.repo.MeasureGetmeas, ctx, acc, params)
}

func (s *service) CallService(ctx context.Context, acc *account.Account, appli int, params string) ([]byte, error) {
	panic("not implemented")
}
