package server

import (
	"awans.org/aft/internal/db"
	"awans.org/aft/internal/oplog"
	"fmt"
	"github.com/NYTimes/gziphandler"
	"log"
	"net/http"
	"time"
)

func Run(dblogPath string) {
	appDB := db.New()
	dbLog, err := oplog.OpenGobLog(dblogPath)
	defer dbLog.Close()
	if err != nil {
		panic(err)
	}
	oplog.DBFromLog(appDB, dbLog)
	appDB = oplog.LoggedDB(dbLog, appDB)
	auditLog := oplog.NewMemLog()

	ops := MakeOps(auditLog)
	router := NewRouter(ops, appDB, auditLog)
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
