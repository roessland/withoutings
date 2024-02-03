package logging

import (
	"context"
	"fmt"
	"github.com/ThreeDotsLabs/watermill"

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

// GetOrCreateLoggerFromContext returns the logger from the context, or creates a new one if it doesn't exist.
// For testing.
func GetOrCreateLoggerFromContext(ctx context.Context) logrus.FieldLogger {
	if log, ok := ctx.Value(ContextKeyLogger).(logrus.FieldLogger); ok {
		return log
	}

	newLog := logrus.New()
	newLog.SetFormatter(&logrus.JSONFormatter{})
	newLog.SetLevel(logrus.DebugLevel)
	return newLog
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

type WatermillLoggerAdapter interface {
	Error(msg string, err error, fields watermill.LogFields)
	Info(msg string, fields watermill.LogFields)
	Debug(msg string, fields watermill.LogFields)
	Trace(msg string, fields watermill.LogFields)
	With(fields watermill.LogFields) watermill.LoggerAdapter
}

type LogrusWatermill struct {
	Logger logrus.FieldLogger
}

func NewLogrusWatermill(logger logrus.FieldLogger) *LogrusWatermill {
	return &LogrusWatermill{Logger: logger}
}

func (l *LogrusWatermill) Error(msg string, err error, fields watermill.LogFields) {
	l.Logger.WithFields(logrus.Fields(fields)).WithError(err).Error(msg)
}

func (l *LogrusWatermill) Info(msg string, fields watermill.LogFields) {
	l.Logger.WithFields(logrus.Fields(fields)).Info(msg)
}

func (l *LogrusWatermill) Debug(msg string, fields watermill.LogFields) {
	l.Logger.WithFields(logrus.Fields(fields)).Debug(msg)
}

func (l *LogrusWatermill) Trace(msg string, fields watermill.LogFields) {
	l.Logger.WithFields(logrus.Fields(fields)).Trace(msg)
}

func (l *LogrusWatermill) With(fields watermill.LogFields) watermill.LoggerAdapter {
	return &LogrusWatermill{Logger: l.Logger.WithFields(logrus.Fields(fields))}
}

var _ WatermillLoggerAdapter = &LogrusWatermill{}
var _ watermill.LoggerAdapter = &LogrusWatermill{}
