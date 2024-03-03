package web

import (
	"github.com/gorilla/mux"
	"github.com/roessland/withoutings/pkg/web/middleware"
	"github.com/roessland/withoutings/pkg/withoutings/app"
)

func AuthMiddlewares(svc *app.App) []mux.MiddlewareFunc {
	return []mux.MiddlewareFunc{
		svc.Sessions.LoadAndSave,
		middleware.Account(svc),
		middleware.FlashMessages(svc),
	}
}

func GlobalMiddlewares(svc *app.App) []mux.MiddlewareFunc {
	return []mux.MiddlewareFunc{
		middleware.UseRemoteAddrFromXForwardedFor(),
		middleware.Logging(svc),
	}
}
