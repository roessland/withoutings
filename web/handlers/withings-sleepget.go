package handlers

//func SleepGetJSON(app *withoutings.Service) http.HandlerFunc {
//	return func(w http.ResponseWriter, r *http.Request) {
//		ctx := r.Context()
//		log := logging.MustGetLoggerFromContext(ctx)
//
//		token := sess.Token()
//
//		var sleepData *withingsapi.SleepGetResponse
//		if token != nil && time.Now().After(token.Expiry) {
//			w.Header().Set("Content-Type", "text/html")
//			w.WriteHeader(200)
//			err = app.Templates.RenderSleepSummaries(w, nil, "Your token is expired. Go refresh it.")
//			if err != nil {
//				log.WithError(err).WithField("event", "error.render.template").Error()
//				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
//				return
//			}
//			return
//		}
//		if token == nil {
//			log.WithError(err).Error("token missing from session")
//			http.Error(w, "token missing from session", http.StatusBadRequest)
//			return
//		}
//
//		authClient := app.Withings.WithAccessToken(token.AccessToken)
//
//		params := withingsapi.NewSleepGetParams()
//		params.Startdate = 1668116907
//		params.Enddate = 1668160107
//		sleepData, err = authClient.SleepGet(ctx, params)
//		if err != nil {
//			log.Error(err)
//			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
//			return
//		}
//
//		w.Header().Set("Content-Type", "application/json")
//		w.WriteHeader(200)
//		w.Write(sleepData.Raw)
//	}
//}
