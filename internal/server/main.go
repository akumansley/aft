package server

import (
	"awans.org/aft/internal/access_log"
	"awans.org/aft/internal/api"
	"awans.org/aft/internal/bus"
	"awans.org/aft/internal/cors"
	"awans.org/aft/internal/db"
	"awans.org/aft/internal/gzip"
	"awans.org/aft/internal/oplog"
	"awans.org/aft/internal/server/lib"
	"fmt"
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
	err = oplog.DBFromLog(appDB, dbLog)
	if err != nil {
		panic(err)
	}
	appDB = oplog.LoggedDB(dbLog, appDB)
	bus := bus.New()
	r := NewRouter()

	modules := []lib.Module{
		gzip.GetModule(),
		cors.GetModule(),
		access_log.GetModule(),
		api.GetModule(appDB, bus),
	}

	for _, mod := range modules {
		r.AddRoutes(mod.ProvideRoutes())
		r.AddMiddleware(mod.ProvideMiddleware())
		bus.RegisterHandlers(mod.ProvideHandlers())
	}

	port := ":8080"
	fmt.Println("Serving on port", port)

	srv := &http.Server{
		Handler:      r,
		Addr:         "localhost:8080",
		WriteTimeout: 1 * time.Second,
		ReadTimeout:  1 * time.Second,
	}

	log.Fatal(srv.ListenAndServeTLS("localhost.pem", "localhost-key.pem"))
}
