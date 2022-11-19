package testctx

import (
	"context"
	"github.com/roessland/withoutings/internal/logging"
	"github.com/sirupsen/logrus"
)

// Context is a context for use in tests.
type Context struct {
	context.Context
}

// New returns a new Context.
func New() Context {
	ctx := context.Background()
	logger := logrus.New()
	ctx = logging.AddLoggerToContext(ctx, logger)
	return Context{ctx}
}
