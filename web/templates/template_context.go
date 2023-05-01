package templates

import (
	"context"
	"github.com/roessland/withoutings/pkg/withoutings/domain/account"
	"github.com/roessland/withoutings/web/flash"
)

type TemplateContext struct {
	Account AccountContext
	Flash   string
}

func extractTemplateContext(ctx context.Context) TemplateContext {
	templateContext := TemplateContext{}

	// Account
	accountCtx := AccountContext{
		IsLoggedIn: false,
	}
	if acc := account.GetAccountFromContext(ctx); acc != nil {
		accountCtx.IsLoggedIn = true
		accountCtx.WithingsUserID = acc.WithingsUserID()
		if acc.CanRefreshAccessToken() {
			accountCtx.AccessTokenState = "stale"
		} else {
			accountCtx.AccessTokenState = "fresh"
		}
	}
	templateContext.Account = accountCtx

	// Flash message
	templateContext.Flash = flash.GetMsgFromContext(ctx)

	return templateContext
}

type AccountContext struct {
	IsLoggedIn       bool
	WithingsUserID   string
	AccessTokenState string
}
