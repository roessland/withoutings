package main

import (
	"github.com/roessland/withoutings/server"
	"github.com/roessland/withoutings/server/serverapp"
	"github.com/roessland/withoutings/worker/workerapp"
	"log"
)

func main() {
	serverApp := serverapp.NewApp()
	srv := server.Configure(serverApp)

	workerApp := workerapp.NewApp()
	go workerApp.Work()

	// 3 - We start up our Client on port 3528
	serverApp.Log.Print("Serving at ", srv.Addr)
	log.Fatal(srv.ListenAndServe())
}
