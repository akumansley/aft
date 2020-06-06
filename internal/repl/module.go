package repl

import (
	"awans.org/aft/internal/bus"
	"awans.org/aft/internal/db"
	"awans.org/aft/internal/server/lib"
	"net/http"
)

type Module struct {
	lib.BlankModule
	db             db.DB
	bus            *bus.EventBus
	dbReadyHandler interface{}
}

func (m Module) ProvideRoutes() []lib.Route {
	return []lib.Route{
		lib.Route{
			Name:    "Repl",
			Pattern: "/views/repl",
			Handler: lib.ErrorHandler(REPLHandler{db: m.db, bus: m.bus}),
		},
	}
}

func GetModule(b *bus.EventBus) lib.Module {
	m := &Module{bus: b}
	m.dbReadyHandler = func(event lib.DatabaseReady) {
		m.db = event.Db
	}
	return m
}

type ErrorResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type apiHandler interface {
	ServeHTTP(http.ResponseWriter, *http.Request) error
}
