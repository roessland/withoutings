package web

import (
	"github.com/gorilla/mux"
	"github.com/roessland/withoutings/pkg/withoutings/app"
	"github.com/roessland/withoutings/web/middleware"
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
