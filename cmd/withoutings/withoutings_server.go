package main

import (
	"context"
	"github.com/roessland/withoutings/internal/services/withoutings"
	"github.com/roessland/withoutings/web"
	"github.com/roessland/withoutings/worker"
	"os"
	"os/signal"
	"runtime/pprof"
)

func withoutingsServer() {
	ctx, cancel := context.WithCancel(context.Background())

	svc, err := withoutings.NewService(ctx)
	if err != nil {
		svc.Log.Fatal(err)
	}

	webserver := web.Server(svc)

	wrk := worker.NewWorker(svc)
	go wrk.Work(ctx)

	go func() {
		svc.Log.Info("Serving at ", webserver.Addr)
		svc.Log.Info(webserver.ListenAndServe())
	}()

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	<-signalChan
	svc.Log.Info("Interrupted. Exiting gracefully. Press Ctrl-C again to exit immediately.")
	go func() {
		<-signalChan
		svc.Log.Info("Interrupted again. Exiting immediately. Running goroutines:")
		_ = pprof.Lookup("goroutine").WriteTo(os.Stdout, 1)
		os.Exit(1)
	}()

	cancel()
	_ = webserver.Close()
}
