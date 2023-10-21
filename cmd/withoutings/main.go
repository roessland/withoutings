package main

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"runtime/pprof"
	"time"

	"github.com/roessland/withoutings/pkg/config"
	"github.com/roessland/withoutings/pkg/web"
	"github.com/roessland/withoutings/pkg/withoutings/app"
	"github.com/roessland/withoutings/worker"
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

	webserver := web.Server(svc)

	// Start worker
	wrk := worker.NewWorker(svc)
	go wrk.Work(ctx)
	svc.Log.WithField("event", "worker.started").Info()

	// Start server
	g.Go(func() error {
		svc.Log.WithField("event", "webserver.started").WithField("addr", webserver.Addr).Info()
		if err := webserver.ListenAndServe(); err != http.ErrServerClosed {
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
	g.Go(func() error {
		<-ctx.Done()
		svc.Log.Info("Interrupted. Attempting graceful shutdown... Wait 10 seconds or press Ctrl-C to exit immediately.")
		deregister()
		signalCh := make(chan os.Signal, 1)
		signal.Notify(signalCh, os.Interrupt)

		select {
		case <-signalCh:
			svc.Log.WithField("event", "graceful-shutdown.interrupted").Error("Exiting immediately...")
		case <-time.After(5 * time.Second):
			svc.Log.WithField("event", "graceful-shutdown.timeout").Error("Exiting immediately...")
		case <-forcefulShutdownCtx.Done():
			return nil
		}

		var buf bytes.Buffer
		_ = pprof.Lookup("goroutine").WriteTo(&buf, 2)
		svc.Log.
			WithField("event", "graceful-shutdown.goroutine-dump").
			WithField("goroutines-remaining", buf.String()).
			Error()
		buf.WriteTo(os.Stdout)

		time.Sleep(100 * time.Millisecond)
		os.Exit(1)
		return nil
	})

	err = g.Wait()
	if err != nil {
		panic(err)
	}

	time.Sleep(15 * time.Second)
	abortForcefulShutdown() // Wasn't needed

	svc.Log.WithField("event", "graceful-shutdown.success").Info()
	time.Sleep(100 * time.Millisecond) // Wait a bit to increase chance of logs being flushed
}
