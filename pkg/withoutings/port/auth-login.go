package port

import (
	"fmt"
	"github.com/roessland/withoutings/pkg/logging"
	"github.com/roessland/withoutings/pkg/withoutings/app"
	"github.com/roessland/withoutings/pkg/withoutings/domain/account"
	"net/http"
)

// LoginPage renders the login page, where users can log in to the application using
// either username/password or Withings OAuth2.
//
// Methods: GET
func LoginPage(svc *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		log := logging.MustGetLoggerFromContext(ctx)

		acc := account.GetFromContext(ctx)

		var errMsg string
		if acc == nil {
			errMsg = ""
		} else {
			errMsg = fmt.Sprintf("You are already logged in as WithingsUserId=%s", acc.WithingsUserID())
		}

		w.WriteHeader(200)
		w.Header().Set("Content-Type", "text/html")
		tmplErr := svc.Templates.RenderLoginPage(ctx, w, errMsg)
		if tmplErr != nil {
			log.WithError(tmplErr).WithField("event", "error.RenderSubscriptionsPage").Error()
			return
		}
	}
}
