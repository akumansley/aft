package server

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"awans.org/aft/internal/access_log"
	"awans.org/aft/internal/api/handlers"
	"awans.org/aft/internal/audit"
	"awans.org/aft/internal/auth"
	"awans.org/aft/internal/bizdatatypes"
	"awans.org/aft/internal/bus"
	"awans.org/aft/internal/catalog"
	"awans.org/aft/internal/cors"
	"awans.org/aft/internal/db"
	"awans.org/aft/internal/explorer"
	"awans.org/aft/internal/gzip"
	"awans.org/aft/internal/oplog"
	"awans.org/aft/internal/rpc"
	"awans.org/aft/internal/server/lib"
	"awans.org/aft/internal/starlark"
)

func Run(dblogPath string, authed bool) {
	bus := bus.New()
	appDB := db.New(bus)

	dbLog, err := oplog.OpenGobLog(dblogPath)
	defer dbLog.Close()
	if err != nil {
		panic(err)
	}

	modules := []lib.Module{
		gzip.GetModule(),
		audit.GetModule(bus, dbLog, oplog.NewMemLog()),
		explorer.GetModule(),
		catalog.GetModule(),
		handlers.GetModule(bus),
		rpc.GetModule(bus),
		bizdatatypes.GetModule(bus),
		starlark.GetModule(),
		access_log.GetModule(),
		auth.GetModule(bus),
		cors.GetModule(),
	}

	for _, mod := range modules {
		bus.RegisterHandlers(mod.ProvideHandlers())
	}

	for _, mod := range modules {
		for _, fl := range mod.ProvideFunctionLoaders() {
			appDB.RegisterRuntime(fl)
		}

		for _, model := range mod.ProvideModels() {
			appDB.AddLiteral(model)
			r := db.RecordForModel(model)
			oplog.Register(r)
		}

		funcs := mod.ProvideFunctions()
		for _, f := range funcs {
			if nf, ok := f.(db.NativeFunctionL); ok {
				appDB.RegisterNativeFunction(nf)
			}
			appDB.AddLiteral(f)
		}
		datatypes := mod.ProvideDatatypes()
		for _, dt := range datatypes {
			appDB.AddLiteral(dt)
		}
		literals := mod.ProvideLiterals()
		for _, lt := range literals {
			appDB.AddLiteral(lt)
		}
	}

	err = oplog.DBFromLog(appDB, dbLog)
	if err != nil {
		panic(err)
	}
	txLogger := oplog.MakeTransactionLogger(dbLog)
	bus.RegisterHandler(txLogger)

	if authed {
		appDB = auth.AuthedDB(appDB)
		bus.RegisterHandler(auth.PostconditionHandler)
	}

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
