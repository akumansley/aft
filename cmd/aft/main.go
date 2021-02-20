package main

import (
	"flag"
	"fmt"
	"os"

	"awans.org/aft/internal/server"
)

var AFT_VERSION = "0.0.5"

func main() {
	dbp := flag.String("db", "", "db log file")
	authed := flag.Bool("authed", true, "enable auth")
	catPort := flag.String("port", "8081", "aft port")
	servePort := flag.String("serve_port", "8080", "app port")
	serveDir := flag.String("serve_dir", "", "app dir")
	tlsKey := flag.String("tls_key", "", "tls key")
	tlsCert := flag.String("tls_cert", "", "tls cert")
	version := flag.Bool("version", false, "print current version")

	flag.Parse()

	if *version {
		fmt.Printf(AFT_VERSION)
		os.Exit(0)
	}

	dblogPath := *dbp
	if dblogPath == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}

	server.Run(
		server.TLSCert(*tlsCert),
		server.TLSKey(*tlsKey),
		server.DBLogPath(dblogPath),
		server.Authed(*authed),
		server.CatalogPort(*catPort),
		server.ServePort(*servePort),
		server.ServeDir(*serveDir),
	)
}
