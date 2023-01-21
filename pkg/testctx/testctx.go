package testctx

import (
	"context"
	"github.com/roessland/withoutings/pkg/logging"
	"github.com/sirupsen/logrus"
)

// Context is a context for use in tests.
type Context struct {
	context.Context
	Logger *logrus.Logger
}

// New returns a new Context.
func New() Context {
	ctx := context.Background()
	logger := logrus.New()
	ctx = logging.AddLoggerToContext(ctx, logger)
	return Context{Context: ctx, Logger: logger}
}
