package middleware

import (
	"bytes"
	"context"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/roessland/withoutings/logging"
	"github.com/roessland/withoutings/server/serverapp"
	"net/http"
)

var ContextKeyRequestID = "requestID"

func Logging(app *serverapp.App) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			log := app.Log

			requestID := uuid.New()
			ctx = context.WithValue(ctx, ContextKeyRequestID, requestID.String())

			log = log.WithField("requestID", requestID)
			log = log.WithField("url", r.URL.String())

			ctx = logging.AddLoggerToContext(ctx, log)

			log.WithField("event", "request.start").
				WithField("headers", r.Header).
				Info("")

			responseRecorder := NewRecordingResponseWriter(w)
			next.ServeHTTP(responseRecorder, r.WithContext(ctx))

			log.WithField("headers", responseRecorder.Header()).
				WithField("response_status", responseRecorder.StatusCode()).
				WithField("event", "request.finish").
				Info()
		})
	}
}

type RecordingResponseWriter struct {
	http.ResponseWriter
	body       bytes.Buffer
	statusCode int
}

func NewRecordingResponseWriter(w http.ResponseWriter) *RecordingResponseWriter {
	rrw := RecordingResponseWriter{}
	rrw.ResponseWriter = w
	return &rrw
}

func (rrw *RecordingResponseWriter) Write(buf []byte) (int, error) {
	rrw.body.Write(buf)
	return rrw.ResponseWriter.Write(buf)
}

func (rrw *RecordingResponseWriter) WriteHeader(statusCode int) {
	rrw.statusCode = statusCode
	rrw.ResponseWriter.WriteHeader(statusCode)
}

func (rrw *RecordingResponseWriter) Body() []byte {
	return rrw.body.Bytes()
}

func (rrw *RecordingResponseWriter) StatusCode() int {
	if rrw.statusCode == 0 {
		return 200
	} else {
		return rrw.statusCode
	}
}
