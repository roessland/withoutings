package withings

import (
	"context"
	"fmt"
	"github.com/roessland/withoutings/pkg/withoutings/domain/account"
	"github.com/roessland/withoutings/pkg/withoutings/domain/withings"
)

type Service struct {
	repo        withings.Repo
	accountRepo account.Repo
}

func NewService(withingsRepo withings.Repo, accountRepo account.Repo) *Service {
	return &Service{
		repo:        withingsRepo,
		accountRepo: accountRepo,
	}
}

// attemptAccessTokenRefresh refreshes the access token for the given account and persists it.
func (s *Service) attemptAccessTokenRefresh(ctx context.Context, acc *account.Account) error {
	// Get new token from Withings API
	token, err := s.repo.RefreshAccessToken(ctx, acc.WithingsRefreshToken())
	refreshSucceeded := err == nil

	// Store updated token in database
	if refreshSucceeded {
		err = s.accountRepo.Update(ctx, acc, func(ctx context.Context, acc *account.Account) (*account.Account, error) {
			err := acc.UpdateWithingsToken(token.AccessToken, token.RefreshToken, token.Expiry, token.Scope)
			if err != nil {
				return nil, fmt.Errorf("UpdateWithingsToken failed: %w", err)
			}
			return acc, nil
		})
		if err != nil {
			return fmt.Errorf("failed to update account with new token: %w", err)
		}
	}
	return nil
}

// NotifyList does a NotifyList request with automated access token renewal if necessary.
// TODO build a generic retry method
func (s *Service) NotifyList(
	ctx context.Context,
	acc *account.Account,
	params withings.NotifyListParams,
) (*withings.NotifyListResponse, error) {
	// Get latest account from database to ensure we have the latest token
	acc, err := s.accountRepo.GetAccountByUUID(ctx, acc.UUID())
	if err != nil {
		return nil, fmt.Errorf("failed to get account in preparation for token refresh: %w", err)
	}

	// First attempt
	resp, err := s.repo.NotifyList(ctx, acc.WithingsAccessToken(), params)
	if err == nil {
		return resp, nil
	}

	// Refresh token
	if err == withings.ErrInvalidToken {
		err = s.attemptAccessTokenRefresh(ctx, acc)
		if err != nil {
			return nil, fmt.Errorf("failed to refresh access token: %w", err)
		}
		acc, err = s.accountRepo.GetAccountByUUID(ctx, acc.UUID())
		if err != nil {
			return nil, fmt.Errorf("failed to get account after token refresh: %w", err)
		}
	} else if err != nil {
		return nil, fmt.Errorf("unexpected error from NotifyList: %w", err)
	}

	// Second attempt using refreshed token
	return s.repo.NotifyList(ctx, acc.WithingsAccessToken(), params)
}
