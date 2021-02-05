package handlers

import (
	"awans.org/aft/internal/api/functions"
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

func (m *Module) ID() db.ID {
	return db.MakeID("25b63f7c-deb0-4402-882d-8f6faed77b05")
}

func (m *Module) Name() string {
	return "api"
}

func (m *Module) Package() string {
	return "awans.org/aft/internal/api/handlers"
}

func (m *Module) ProvideRoutes() []lib.Route {
	return []lib.Route{
		lib.Route{
			Name:    "FindOne",
			Pattern: "/api/{modelName}.findOne",
			Handler: lib.ErrorHandler(FindOneHandler{db: m.db, bus: m.b}),
		},
		lib.Route{
			Name:    "FindMany",
			Pattern: "/api/{modelName}.findMany",
			Handler: lib.ErrorHandler(FindManyHandler{db: m.db, bus: m.b}),
		},
		lib.Route{
			Name:    "Create",
			Pattern: "/api/{modelName}.create",
			Handler: lib.ErrorHandler(CreateHandler{DB: m.db, Bus: m.b}),
		},
		lib.Route{
			Name:    "Delete",
			Pattern: "/api/{modelName}.delete",
			Handler: lib.ErrorHandler(DeleteHandler{db: m.db, bus: m.b}),
		},
		lib.Route{
			Name:    "DeleteMany",
			Pattern: "/api/{modelName}.deleteMany",
			Handler: lib.ErrorHandler(DeleteManyHandler{db: m.db, bus: m.b}),
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
		lib.Route{
			Name:    "Count",
			Pattern: "/api/{modelName}.count",
			Handler: lib.ErrorHandler(CountHandler{db: m.db, bus: m.b}),
		},
		lib.Route{
			Name:    "Upsert",
			Pattern: "/api/{modelName}.upsert",
			Handler: lib.ErrorHandler(UpsertHandler{db: m.db, bus: m.b}),
		},
	}
}

func GetModule(b *bus.EventBus) lib.Module {
	m := &Module{b: b}
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

func (m *Module) ProvideFunctions() []db.FunctionL {
	return []db.FunctionL{
		functions.FindOneFunc,
		functions.FindManyFunc,
		functions.CountFunc,
		functions.DeleteFunc,
		functions.DeleteManyFunc,
		functions.UpdateFunc,
		functions.UpdateManyFunc,
		functions.CreateFunc,
		functions.UpsertFunc,
	}
}
