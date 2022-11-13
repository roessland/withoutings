package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		printUsage()
	}
	switch os.Args[1] {
	case "server":
		withoutingsServer()
	case "migrate":
		withoutingsMigrate()
	default:
		printUsage()
	}
}

func printUsage() {
	fmt.Fprintf(os.Stderr, "usage: withoutings migrate|server")
	os.Exit(1)
}
