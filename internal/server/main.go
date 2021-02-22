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
		handlers.GetModule(),
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
		for _, fl := range mod.ProvideFunctionLoaders() {
			appDB.RegisterRuntime(fl)
			appDB.AddLiteral(fl.ProvideModel())
		}

		funcs := mod.ProvideFunctions()
		for _, f := range funcs {
			if nf, ok := f.(db.NativeFunctionLiteral); ok {
				appDB.RegisterNativeFunction(nf)
			}
			appDB.AddLiteral(f)
		}

		datatypes := mod.ProvideDatatypes()
		for _, dt := range datatypes {
			appDB.AddLiteral(dt)
		}

		for _, iface := range mod.ProvideInterfaces() {
			appDB.AddLiteral(iface)
		}

		literals := mod.ProvideLiterals()
		for _, lt := range literals {
			appDB.AddLiteral(lt)
		}

		appDB.AddLiteral(lib.ToLiteral(mod))
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
		spaHandler := &spaHandler{
			Dir: http.Dir(c.ServeDir),
		}
		go func() {
			serve("client", c, modules, spaHandler, servePort)
		}()
	}

	port := c.CatalogPort
	handler := &spaHandler{Dir: serverCatalog.Dir}
	serve("dev", c, modules, handler, port)
}

func serve(name string, c *config, modules []lib.Module, h http.Handler, port string) {
	useTLS := c.TLSKey != ""
	scheme := "http://"
	if useTLS {
		scheme = "https://"
	}
	path := fmt.Sprintf("localhost:%v", port)

	r := NewRouter(h)
	for _, mod := range modules {
		r.AddRoutes(mod.ProvideRoutes())
		r.AddMiddleware(mod.ProvideMiddleware())
	}

	server := &http.Server{
		Handler:      r,
		Addr:         path,
		WriteTimeout: 1 * time.Second,
		ReadTimeout:  1 * time.Second,
	}

	fmt.Printf("serving %v on %v%v\n", name, scheme, path)
	if useTLS {
		log.Fatal(server.ListenAndServeTLS(c.TLSCert, c.TLSKey))
	} else {
		log.Fatal(server.ListenAndServe())
	}
}
