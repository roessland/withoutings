package main

import (
	"github.com/roessland/withoutings/web"
	"github.com/roessland/withoutings/web/webapp"
	"github.com/roessland/withoutings/worker/workerapp"
	"log"
)

func main() {
	app := app.NewApp()
	srv := web.Configure(app)

	workerApp := workerapp.NewApp()
	go workerApp.Work()

	app.Log.Print("Serving at ", srv.Addr)
	log.Fatal(srv.ListenAndServe())
}
