package withings

import (
	"context"
)

//go:generate mockery --name Repo --filename domain_withingsrepo_mock.go
type Repo interface {
	AuthCodeURL(nonce string) string
	GetAccessToken(ctx context.Context, authCode string) (*Token, error)
	RefreshAccessToken(ctx context.Context, refreshToken string) (*Token, error)
	NotifySubscribe(ctx context.Context, accessToken string, params NotifySubscribeParams) (*NotifySubscribeResponse, error)
	NotifyList(ctx context.Context, accessToken string, params NotifyListParams) (*NotifyListResponse, error)
	SleepGetsummary(ctx context.Context, accessToken string, params SleepGetSummaryParams) (*SleepGetsummaryResponse, error)
	SleepGet(ctx context.Context, accessToken string, params SleepGetParams) (*SleepGetResponse, error)
}
