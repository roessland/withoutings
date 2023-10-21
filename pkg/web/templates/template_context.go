package templates

import (
	"context"

	"github.com/roessland/withoutings/pkg/web/flash"
	"github.com/roessland/withoutings/pkg/withoutings/domain/account"
)

type TemplateContext struct {
	Flash   string
	Account AccountContext
}

func extractTemplateContext(ctx context.Context) TemplateContext {
	templateContext := TemplateContext{}

	// Account
	accountCtx := AccountContext{
		IsLoggedIn: false,
	}
	if acc := account.GetFromContext(ctx); acc != nil {
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
	WithingsUserID   string
	AccessTokenState string
	IsLoggedIn       bool
}
