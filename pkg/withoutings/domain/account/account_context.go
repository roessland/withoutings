package account

import "context"

type contextKey struct{}

var contextKeyAccount contextKey = struct{}{}

func GetFromContext(ctx context.Context) *Account {
	acc, ok := ctx.Value(contextKeyAccount).(*Account)
	if !ok {
		return nil
	}
	return acc
}

func AddToContext(ctx context.Context, account *Account) context.Context {
	return context.WithValue(ctx, contextKeyAccount, account)
}
