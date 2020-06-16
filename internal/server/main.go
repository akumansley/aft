package server

import (
	"awans.org/aft/internal/access_log"
	"awans.org/aft/internal/api"
	"awans.org/aft/internal/auth"
	"awans.org/aft/internal/bizdatatypes"
	"awans.org/aft/internal/bus"
	"awans.org/aft/internal/cors"
	"awans.org/aft/internal/db"
	"awans.org/aft/internal/gzip"
	"awans.org/aft/internal/oplog"
	"awans.org/aft/internal/repl"
	"awans.org/aft/internal/runtime"
	"awans.org/aft/internal/server/lib"
	"fmt"
	"log"
	"net/http"
	"time"
)

func Run(dblogPath string) {
	bus := bus.New()
	appDB := db.New(&runtime.Executor{})

	modules := []lib.Module{
		gzip.GetModule(),
		cors.GetModule(),
		access_log.GetModule(),
		api.GetModule(bus),
		auth.GetModule(bus),
		bizdatatypes.GetModule(bus),
		repl.GetModule(bus),
	}

	for _, mod := range modules {
		bus.RegisterHandlers(mod.ProvideHandlers())
	}

	tx := appDB.NewRWTx()
	for _, mod := range modules {
		for _, model := range mod.ProvideModels() {
			tx.SaveModel(model)
		}
		for _, record := range mod.ProvideRecords() {
			tx.Insert(record)
		}
	}
	tx.Commit()

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

	bus.Publish(lib.DatabaseReady{Db: appDB})

	r := NewRouter()

	for _, mod := range modules {
		r.AddRoutes(mod.ProvideRoutes())
		r.AddMiddleware(mod.ProvideMiddleware())
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
