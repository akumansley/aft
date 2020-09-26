package main

import (
	"flag"
	"os"

	"awans.org/aft/internal/server"
)

func main() {
	dbp := flag.String("db", "", "db log file")
	authed := flag.Bool("authed", true, "enable auth")
	flag.Parse()
	dblogPath := *dbp
	if dblogPath == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}
	server.Run(dblogPath, *authed)
}
