package rpc

import (
	"awans.org/aft/internal/bus"
	"awans.org/aft/internal/db"
	"awans.org/aft/internal/server/lib"
)

type Module struct {
	lib.BlankModule
	db             db.DB
	bus            *bus.EventBus
	dbReadyHandler interface{}
	authed         bool
}

func (m *Module) ID() db.ID {
	return db.MakeID("d7bdfdc9-9bbb-4083-b46f-868ad1552a4c")
}

func (m *Module) Name() string {
	return "rpc"
}

func (m *Module) Package() string {
	return "awans.org/aft/internal/rpc"
}

func (m Module) ProvideRoutes() []lib.Route {
	return []lib.Route{
		lib.Route{
			Name:    "RPC",
			Pattern: "/rpc/{name}",
			Handler: lib.ErrorHandler(RPCHandler{db: m.db, bus: m.bus, authed: m.authed}),
		},
	}
}

func GetModule(b *bus.EventBus, authed bool) lib.Module {
	m := &Module{bus: b, authed: authed}
	m.dbReadyHandler = func(event lib.DatabaseReady) {
		m.db = event.DB
	}
	return m
}

func (m *Module) ProvideHandlers() []interface{} {
	return []interface{}{
		m.dbReadyHandler,
	}
}
