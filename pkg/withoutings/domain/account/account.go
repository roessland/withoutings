package account

import (
	"errors"
	"github.com/google/uuid"
	"time"
)

type Account struct {
	accountUUID               uuid.UUID
	withingsUserID            string
	withingsAccessToken       string
	withingsRefreshToken      string
	withingsAccessTokenExpiry time.Time
	withingsScopes            string
}

func NewAccount(
	accountUUID uuid.UUID,
	withingsUserID string,
	withingsAccessToken string,
	withingsRefreshToken string,
	withingsAccessTokenExpiry time.Time,
	withingsScopes string,
) (*Account, error) {
	if accountUUID == uuid.Nil {
		return nil, errors.New("empty account UUID")
	}
	if withingsUserID == "" {
		return nil, errors.New("empty withings user ID")
	}
	if withingsAccessToken == "" {
		return nil, errors.New("empty access token")
	}
	if withingsRefreshToken == "" {
		return nil, errors.New("empty refresh token")
	}
	if withingsAccessTokenExpiry.IsZero() {
		return nil, errors.New("zero expiry time")
	}
	if withingsScopes == "" {
		return nil, errors.New("empty scopes")
	}

	return &Account{
		accountUUID:               accountUUID,
		withingsUserID:            withingsUserID,
		withingsAccessToken:       withingsAccessToken,
		withingsRefreshToken:      withingsRefreshToken,
		withingsAccessTokenExpiry: withingsAccessTokenExpiry,
		withingsScopes:            withingsScopes,
	}, nil
}

// NewAccountOrPanic is a convenience function for tests.
func NewAccountOrPanic(
	accountUUID uuid.UUID,
	withingsUserID string,
	withingsAccessToken string,
	withingsRefreshToken string,
	withingsAccessTokenExpiry time.Time,
	withingsScopes string,
) *Account {
	acc, err := NewAccount(
		accountUUID,
		withingsUserID,
		withingsAccessToken,
		withingsRefreshToken,
		withingsAccessTokenExpiry,
		withingsScopes,
	)
	if err != nil {
		panic(err)
	}
	return acc
}

func (acc *Account) UUID() uuid.UUID {
	return acc.accountUUID
}

func (acc *Account) WithingsUserID() string {
	return acc.withingsUserID
}

func (acc *Account) WithingsAccessToken() string {
	return acc.withingsAccessToken
}

func (acc *Account) WithingsRefreshToken() string {
	return acc.withingsRefreshToken
}

func (acc *Account) WithingsAccessTokenExpiry() time.Time {
	return acc.withingsAccessTokenExpiry
}

func (acc *Account) WithingsScopes() string {
	return acc.withingsScopes
}

func (acc *Account) CanRefreshAccessToken() bool {
	return time.Now().After(acc.withingsAccessTokenExpiry)
}

func (acc *Account) UpdateWithingsToken(accessToken string, refreshToken string, expiry time.Time, scopes string) error {
	if accessToken == "" {
		return errors.New("empty access token")
	}
	if refreshToken == "" {
		return errors.New("empty refresh token")
	}
	if expiry.Before(time.Now()) {
		return errors.New("already expired")
	}
	if scopes == "" {
		return errors.New("empty scopes")
	}
	acc.withingsAccessToken = accessToken
	acc.withingsRefreshToken = refreshToken
	acc.withingsAccessTokenExpiry = expiry
	acc.withingsScopes = scopes
	return nil
}
