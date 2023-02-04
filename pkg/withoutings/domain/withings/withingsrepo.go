package withings

import (
	"context"
	"fmt"
)

type TodoError struct {
	TodoProperty int64
}

func (e TodoError) Error() string {
	return fmt.Sprintf("todo with ID %d not found", e.TodoProperty)
}

type Repo interface {
	AuthCodeURL(nonce string) string
	GetAccessToken(ctx context.Context, authCode string) (*Token, error)
	RefreshAccessToken(ctx context.Context, refreshToken string) (*Token, error)
	NotifySubscribe(ctx context.Context, accessToken string, params NotifySubscribeParams) (*NotifySubscribeResponse, error)
	SleepGetsummary(ctx context.Context, accessToken string, params SleepGetSummaryParams) (*SleepGetsummaryResponse, error)
	SleepGet(ctx context.Context, accessToken string, params SleepGetParams) (*SleepGetResponse, error)
}
