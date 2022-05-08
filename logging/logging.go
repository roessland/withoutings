package logging

import (
	"context"
	"github.com/sirupsen/logrus"
)

var ContextKeyLogger = "logger"

func MustGetLoggerFromContext(ctx context.Context) logrus.FieldLogger {
	log, ok := ctx.Value(ContextKeyLogger).(logrus.FieldLogger)
	if !ok {
		panic("no logger on context")
	}
	return log
}

func AddLoggerToContext(ctx context.Context, log logrus.FieldLogger) context.Context {
	return context.WithValue(ctx, ContextKeyLogger, log)
}
