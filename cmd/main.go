package main

import (
	"github.com/roessland/withoutings/server"
	"github.com/roessland/withoutings/server/app"
	"log"
)

func main() {
	app := app.NewApp()
	srv := server.Configure(app)

	// 3 - We start up our Client on port 9094
	app.Log.Print("Serving at ", srv.Addr)
	log.Fatal(srv.ListenAndServe())
}
