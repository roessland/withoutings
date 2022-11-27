package withingsapi

type Client interface {
	NotifySubscribe(ctx context.Context, params NotifySubscribeParams) error {
}
