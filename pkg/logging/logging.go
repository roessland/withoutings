package logging

import (
	"context"
	"fmt"

	"github.com/sirupsen/logrus"
)

type symbol string

var ContextKeyLogger = symbol("logger")

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

func NewLogger(logFormat string) *logrus.Logger {
	switch logFormat {
	case "":
		fallthrough
	case "json":
		log := logrus.New()
		log.SetFormatter(&logrus.JSONFormatter{})
		return log
	case "text":
		log := logrus.New()
		log.SetFormatter(&logrus.TextFormatter{})
		return log
	default:
		panic(fmt.Sprintf("unknown log format: %s", logFormat))
	}
}
