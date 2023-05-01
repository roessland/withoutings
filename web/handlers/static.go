package handlers

import (
	"github.com/roessland/withoutings/pkg/withoutings/app"
	"github.com/roessland/withoutings/web/static"
	"net/http"
	"strings"
)

func Static(svc *app.App) http.HandlerFunc {
	h := http.FileServer(http.FS(static.FS))
	return func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, ".css") {
			w.Header().Set("Content-Type", "text/css")
		}
		h.ServeHTTP(w, r)
	}
}
