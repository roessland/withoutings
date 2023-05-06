package web

import (
	"github.com/gorilla/mux"
	"github.com/roessland/withoutings/pkg/web/middleware"
	"github.com/roessland/withoutings/pkg/withoutings/app"
)

func Middleware(svc *app.App) []mux.MiddlewareFunc {
	return []mux.MiddlewareFunc{
		middleware.UseRemoteAddrFromXForwardedFor(),
		middleware.Logging(svc),
		svc.Sessions.LoadAndSave,
		middleware.Account(svc),
		middleware.FlashMessages(svc),
	}
}
