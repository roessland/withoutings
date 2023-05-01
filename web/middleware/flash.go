package middleware

import (
	"github.com/gorilla/mux"
	"github.com/roessland/withoutings/pkg/withoutings/app"
	"github.com/roessland/withoutings/web/flash"
	"net/http"
)

func FlashMessages(svc *app.App) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			flashMessage := svc.Flash.PopMsg(ctx)
			if flashMessage != "" {
				ctx = flash.AddMsgToContext(ctx, flashMessage)
			}

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
