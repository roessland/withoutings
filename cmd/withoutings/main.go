package main

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/roessland/withoutings/pkg/logging"
	"github.com/roessland/withoutings/pkg/worker"
	"net/http"
	"os"
	"os/signal"
	"runtime/pprof"
	"time"

	"github.com/roessland/withoutings/pkg/config"
	"github.com/roessland/withoutings/pkg/web"
	"github.com/roessland/withoutings/pkg/withoutings/app"
	"golang.org/x/sync/errgroup"
)

func withoutingsServer() {
	ctx, deregister := signal.NotifyContext(context.Background(), os.Interrupt)
	defer deregister()

	cfg, err := config.LoadFromEnv()
	if err != nil {
		panic(fmt.Sprintf("load config: %s", err))
	}

	g, ctx := errgroup.WithContext(ctx)

	svc := app.NewApplication(ctx, cfg)
	ctx = logging.AddLoggerToContext(ctx, svc.Log)

	webserver := web.Server(svc)

	// Start worker
	wrk := worker.NewWorker(svc)
	go wrk.Work(ctx)
	svc.Log.WithField("event", "info.worker.started").Info()

	// Start server
	g.Go(func() error {
		svc.Log.WithField("event", "info.webserver.started").WithField("addr", webserver.Addr).Info()
		if err := webserver.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			return err
		}
		return nil
	})

	// On shutdown, close the server
	g.Go(func() error {
		<-ctx.Done()
		return webserver.Shutdown(ctx)
	})

	// On shutdown, do forceful exit after another interrupt, or after 10 seconds
	forcefulShutdownCtx, abortForcefulShutdown := context.WithCancel(context.WithoutCancel(ctx))
	go (func() error {
		<-ctx.Done()

		svc.Log.
			WithField("event", "graceful-shutdown.initiated").
			Info("Interrupted. Attempting graceful shutdown... Wait 10 seconds or press Ctrl-C to exit immediately.")

		deregister()
		signalCh := make(chan os.Signal, 1)
		signal.Notify(signalCh, os.Interrupt)

		select {
		case <-signalCh:
			svc.Log.
				WithField("event", "graceful-shutdown.interrupted").Error()
		case <-time.After(5 * time.Second):
			svc.Log.
				WithField("event", "graceful-shutdown.timeout").Error()
		case <-forcefulShutdownCtx.Done():
			return nil
		}

		var buf bytes.Buffer
		_ = pprof.Lookup("goroutine").WriteTo(&buf, 2)
		svc.Log.
			WithField("event", "graceful-shutdown.goroutine-dump").
			WithField("goroutines-remaining", buf.String()).Error()

		// Useful for development
		if svc.Config.LogFormat == "text" {
			buf.WriteTo(os.Stdout)
		}

		time.Sleep(100 * time.Millisecond)
		os.Exit(1)
		return nil
	})()

	err = g.Wait()
	if err != nil {
		panic(err)
	}

	abortForcefulShutdown() // Wasn't needed after all

	svc.Log.WithField("event", "info.graceful-shutdown.success").Info()
	time.Sleep(100 * time.Millisecond) // Wait a bit to increase chance of logs being flushed
}
