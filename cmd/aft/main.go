package main

import (
	"flag"
	"os"

	"awans.org/aft/internal/server"
)

func main() {
	dbp := flag.String("db", "", "db log file")
	authed := flag.Bool("authed", true, "enable auth")
	catPort := flag.String("port", "8081", "aft port")
	servePort := flag.String("serve_port", "8080", "app port")
	serveDir := flag.String("serve_dir", "", "app dir")

	flag.Parse()
	
	dblogPath := *dbp
	if dblogPath == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}
	server.Run(
		server.DBLogPath(dblogPath),
		server.Authed(*authed),
		server.CatalogPort(*catPort),
		server.ServePort(*servePort),
		server.ServeDir(*serveDir),
	)
}
