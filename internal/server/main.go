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
	authRPCs "awans.org/aft/internal/auth/rpcs"
	"awans.org/aft/internal/bus"
	"awans.org/aft/internal/catalog"
	"awans.org/aft/internal/csrf"
	"awans.org/aft/internal/db"
	"awans.org/aft/internal/explorer"
	"awans.org/aft/internal/gzip"
	"awans.org/aft/internal/oplog"
	"awans.org/aft/internal/phone"
	"awans.org/aft/internal/rpc"
	serverCatalog "awans.org/aft/internal/server/catalog"
	"awans.org/aft/internal/server/lib"
	"awans.org/aft/internal/starlark"
	"awans.org/aft/internal/url"
)

func Run(options ...Option) {
	c := newConfig()
	for _, opt := range options {
		opt(c)
	}

	bus := bus.New()
	appDB := db.New(bus)

	dbLogStore, err := oplog.OpenDiskLog(c.DBLogPath)
	if err != nil {
		panic(err)
	}
	dbLog := db.DBOpLog(appDB.Builder(), dbLogStore)
	defer dbLog.Close()

	modules := []lib.Module{
		gzip.GetModule(),
		handlers.GetModule(bus),
		auth.GetModule(bus),
		rpc.GetModule(bus, c.Authed),
		audit.GetModule(bus, dbLog),
		starlark.GetModule(),
		url.GetModule(),
		phone.GetModule(),
		catalog.GetModule(),
		explorer.GetModule(),
		access_log.GetModule(),
		csrf.GetModule(),
		authRPCs.GetModule(),
	}

	for _, mod := range modules {
		bus.RegisterHandlers(mod.ProvideHandlers())
	}

	for _, mod := range modules {
		rwtx := appDB.NewRWTx()
		for _, fl := range mod.ProvideFunctionLoaders() {
			appDB.RegisterRuntime(fl)
			appDB.AddLiteral(rwtx, fl.ProvideModel())
		}

		for _, iface := range mod.ProvideInterfaces() {
			appDB.AddLiteral(rwtx, iface)
		}

		funcs := mod.ProvideFunctions()
		for _, f := range funcs {
			if nf, ok := f.(db.NativeFunctionLiteral); ok {
				appDB.RegisterNativeFunction(nf)
			}
			appDB.AddLiteral(rwtx, f)
		}

		datatypes := mod.ProvideDatatypes()
		for _, dt := range datatypes {
			appDB.AddLiteral(rwtx, dt)
		}
		rwtx.Commit()

		//TODO make AddLiteral just commit
		rwtx = appDB.NewRWTx()
		literals := mod.ProvideLiterals()
		for _, lt := range literals {
			appDB.AddLiteral(rwtx, lt)
		}
		rwtx.Commit()

		rwtx = appDB.NewRWTx()
		appDB.AddLiteral(rwtx, lib.ToLiteral(mod))
		rwtx.Commit()
	}

	err = db.DBFromLog(appDB, dbLog)
	if err != nil {
		panic(err)
	}

	// we've replayed the database; ready it for new tx
	txLogger := db.MakeTransactionLogger(dbLog)
	bus.RegisterHandler(txLogger)

	if c.Authed {
		appDB = auth.AuthedDB(appDB)
		bus.RegisterHandler(auth.PostconditionHandler)
	}

	bus.Publish(lib.DatabaseReady{DB: appDB})

	if c.ServeDir != "" {
		servePort := c.ServePort
		servePath := fmt.Sprintf("localhost:%v", servePort)
		fmt.Printf("Serving client on http://%v\n", servePath)

		spaHandler := &spaHandler{
			Dir: http.Dir(c.ServeDir),
		}
		r := NewRouter(spaHandler)

		for _, mod := range modules {
			r.AddRoutes(mod.ProvideRoutes())
			r.AddMiddleware(mod.ProvideMiddleware())
		}
		spaSrv := &http.Server{
			Handler:      r,
			Addr:         servePath,
			WriteTimeout: 1 * time.Second,
			ReadTimeout:  1 * time.Second,
		}
		go func() {
			log.Fatal(spaSrv.ListenAndServe())
		}()
	}

	port := c.CatalogPort
	path := fmt.Sprintf("localhost:%v", port)
	fmt.Printf("Serving dev on http://%v\n", path)

	r := NewRouter(&spaHandler{Dir: serverCatalog.Dir})

	for _, mod := range modules {
		r.AddRoutes(mod.ProvideRoutes())
		r.AddMiddleware(mod.ProvideMiddleware())
	}

	srv := &http.Server{
		Handler:      r,
		Addr:         path,
		WriteTimeout: 1 * time.Second,
		ReadTimeout:  1 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}
