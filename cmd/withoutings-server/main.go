package main

import (
	"context"
	"github.com/roessland/withoutings/app/webapp"
	"github.com/roessland/withoutings/app/workerapp"
	"github.com/roessland/withoutings/web"
	"os"
	"os/signal"
	"runtime/pprof"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())

	app := webapp.NewApp(ctx)
	server := web.Configure(app)

	worker := workerapp.NewApp()
	go worker.Work(ctx)

	go func() {
		app.Log.Info("Serving at ", server.Addr)
		app.Log.Info(server.ListenAndServe())
	}()

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	<-signalChan
	app.Log.Info("Interrupted. Exiting gracefully. Press Ctrl-C again to exit immediately.")
	go func() {
		<-signalChan
		app.Log.Info("Interrupted again. Exiting immediately. Running goroutines:")
		_ = pprof.Lookup("goroutine").WriteTo(os.Stdout, 1)
		os.Exit(1)
	}()

	cancel()
	_ = server.Close()
}
