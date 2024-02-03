package testctx

import (
	"context"
	"github.com/roessland/withoutings/pkg/logging"
	"github.com/sirupsen/logrus"
)

// Context is a context for use in tests.
type Context struct {
	context.Context
	Logger        *logrus.Logger
	CancelContext context.CancelFunc
}

// New returns a new Context.
func New() Context {
	ctx, cancel := context.WithCancel(context.Background())
	log := logging.NewLogger("json")
	log.SetLevel(logrus.DebugLevel)
	ctx = logging.AddLoggerToContext(ctx, log)
	return Context{Context: ctx, Logger: log, CancelContext: cancel}
}
