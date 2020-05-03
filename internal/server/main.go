package server

import (
	"awans.org/aft/internal/oplog"
	"awans.org/aft/internal/server/db"
	"fmt"
	"github.com/NYTimes/gziphandler"
	"log"
	"net/http"
	"time"
)

func Run() {
	appDB := db.New()
	dbLog := oplog.NewMemLog()
	auditLog := oplog.NewMemLog()
	appDB = oplog.LoggedDB(dbLog, appDB)

	ops := MakeOps(appDB, auditLog)
	router := NewRouter(ops, auditLog)
	port := ":8080"
	fmt.Println("Serving on port", port)
	gzr := gziphandler.GzipHandler(router)

	srv := &http.Server{
		Handler:      gzr,
		Addr:         "localhost:8080",
		WriteTimeout: 1 * time.Second,
		ReadTimeout:  1 * time.Second,
	}

	log.Fatal(srv.ListenAndServeTLS("localhost.pem", "localhost-key.pem"))
}
