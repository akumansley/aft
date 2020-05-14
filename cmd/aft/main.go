package main

import (
	"awans.org/aft/internal/server"
	"flag"
	"os"
)

func main() {
	dbp := flag.String("db", "", "db log file")
	flag.Parse()
	dblogPath := *dbp
	if dblogPath == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}
	server.Run(dblogPath)
}
