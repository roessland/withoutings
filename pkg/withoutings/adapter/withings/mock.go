package withings

import (
	"context"
	"github.com/roessland/withoutings/pkg/withoutings/domain/withings"
)

type MockClient struct{}

func NewMockClient() *MockClient {
	return &MockClient{}
}

func (m MockClient) AuthCodeURL(nonce string) string {
	//TODO implement me
	panic("implement me")
}

func (m MockClient) GetAccessToken(ctx context.Context, authCode string) (*withings.Token, error) {
	return &withings.Token{}, nil
}

func (m MockClient) RefreshAccessToken(ctx context.Context, refreshToken string) (*withings.Token, error) {
	//TODO implement me
	panic("implement me")
}

func (m MockClient) NotifySubscribe(ctx context.Context, accessToken string, params withings.NotifySubscribeParams) (*withings.NotifySubscribeResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (m MockClient) SleepGetsummary(ctx context.Context, accessToken string, params withings.SleepGetSummaryParams) (*withings.SleepGetsummaryResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (m MockClient) SleepGet(ctx context.Context, accessToken string, params withings.SleepGetParams) (*withings.SleepGetResponse, error) {
	//TODO implement me
	panic("implement me")
}
