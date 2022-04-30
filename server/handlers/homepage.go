package handlers

import (
	"github.com/roessland/withoutings/server/app"
	"net/http"
)

func HomePage(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authCodeURL := app.WithingsClient.OAuth2Config.AuthCodeURL("xyfdsfdsz")
		http.Redirect(w, r, authCodeURL, http.StatusFound)
	}
}
