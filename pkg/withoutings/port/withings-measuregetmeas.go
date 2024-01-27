package port

import (
	"fmt"
	"github.com/roessland/withoutings/pkg/logging"
	"github.com/roessland/withoutings/pkg/withoutings/app"
	"github.com/roessland/withoutings/pkg/withoutings/domain/account"
	"github.com/roessland/withoutings/pkg/withoutings/domain/withings"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
)

// MeasureGetmeas renders the getmeas page, showing getmeas responses for
// arbitrary parameters.
func MeasureGetmeas(svc *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		log := logging.MustGetLoggerFromContext(ctx)

		acc := account.GetFromContext(ctx)
		if acc == nil {
			w.Header().Set("Content-Type", "text/html")
			w.WriteHeader(403)
			err := svc.Templates.RenderMeasureGetmeas(ctx, w, "", "You must be logged in to view this page.")
			if err != nil {
				log.WithError(err).WithField("event", "error.render.template").Error()
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				return
			}
			return
		}

		query := r.URL.Query()
		query["userid"] = []string{acc.WithingsUserID()}
		query["action"] = []string{"getmeas"}

		if err := validateQuery(query); err != nil {
			w.Header().Set("Content-Type", "text/html")
			w.WriteHeader(400)
			err = svc.Templates.RenderMeasureGetmeas(ctx, w, "", err.Error())
			if err != nil {
				log.WithError(err).WithField("event", "error.render.template").Error()
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				return
			}
			return
		}

		getmeasBody, err := svc.WithingsSvc.MeasureGetmeas(ctx, acc, withings.MeasureGetmeasParams(query.Encode()))
		if err != nil {
			log.Error(err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(200)
		err = svc.Templates.RenderMeasureGetmeas(ctx, w, string(getmeasBody.Raw), "")
		if err != nil {
			log.WithError(err).WithField("event", "error.render.template").Error()
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
	}

}

func validateQuery(query url.Values) error {
	for k, v := range query {
		if len(v) > 1 {
			return fmt.Errorf("duplicate param")
		}
		switch k {

		case "action":
			if v[0] != "getmeas" {
				return fmt.Errorf("invalid action")
			}

		case "meastype":
			if len(v) == 0 {
				continue
			}
			if _, ok := strconv.Atoi(v[0]); ok != nil {
				return fmt.Errorf("invalid meastype")
			}

		case "meastypes":
			if len(v) == 0 {
				continue
			}
			re := regexp.MustCompile(`^(\d+,)*\d+$`)
			if !re.MatchString(v[0]) {
				return fmt.Errorf("invalid meastypes")
			}

		case "category":
			if len(v) == 0 {
				continue
			}
			if _, ok := strconv.Atoi(v[0]); ok != nil {
				return fmt.Errorf("invalid category")
			}

		case "startdate":
			if len(v) == 0 {
				continue
			}
			if _, ok := strconv.Atoi(v[0]); ok != nil {
				return fmt.Errorf("invalid startdate")
			}

		case "enddate":
			if len(v) == 0 {
				continue
			}
			if _, ok := strconv.Atoi(v[0]); ok != nil {
				return fmt.Errorf("invalid enddate")
			}

		case "lastupdate":
			if len(v) == 0 {
				continue
			}
			if _, ok := strconv.Atoi(v[0]); ok != nil {
				return fmt.Errorf("invalid lastupdate")
			}

		case "offset":
			if len(v) == 0 {
				continue
			}
			if _, ok := strconv.Atoi(v[0]); ok != nil {
				return fmt.Errorf("invalid offset")
			}
		}
	}

	return nil
}
