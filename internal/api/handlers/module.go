package handlers

import (
	"awans.org/aft/internal/bus"
	"awans.org/aft/internal/db"
	"awans.org/aft/internal/server/lib"
)

type Module struct {
	lib.BlankModule
	db             db.DB
	b              *bus.EventBus
	dbReadyHandler interface{}
}

func (m *Module) ProvideRoutes() []lib.Route {
	return []lib.Route{
		lib.Route{
			Name:    "FindMany",
			Pattern: "/api/{modelName}.findMany",
			Handler: lib.ErrorHandler(FindManyHandler{db: m.db, bus: m.b}),
		},
		lib.Route{
			Name:    "FindOne",
			Pattern: "/api/{modelName}.findOne",
			Handler: lib.ErrorHandler(FindOneHandler{db: m.db, bus: m.b}),
		},
		lib.Route{
			Name:    "Create",
			Pattern: "/api/{modelName}.create",
			Handler: lib.ErrorHandler(CreateHandler{DB: m.db, Bus: m.b}),
		},
		lib.Route{
			Name:    "Update",
			Pattern: "/api/{modelName}.update",
			Handler: lib.ErrorHandler(UpdateHandler{db: m.db, bus: m.b}),
		},
		lib.Route{
			Name:    "UpdateMany",
			Pattern: "/api/{modelName}.updateMany",
			Handler: lib.ErrorHandler(UpdateManyHandler{db: m.db, bus: m.b}),
		},
	}
}

func GetModule(b *bus.EventBus) lib.Module {
	m := &Module{b: b}
	m.dbReadyHandler = func(event lib.DatabaseReady) {
		m.db = event.Db
	}
	return m
}

func (m *Module) ProvideHandlers() []interface{} {
	return []interface{}{
		m.dbReadyHandler,
	}
}
