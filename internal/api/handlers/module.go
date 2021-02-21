package handlers

import (
	"awans.org/aft/internal/api/functions"
	"awans.org/aft/internal/db"
	"awans.org/aft/internal/server/lib"
)

type Module struct {
	lib.BlankModule
	db             db.DB
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
			Name:    "API",
			Pattern: "/api/{modelName}.{methodName}",
			Handler: lib.ErrorHandler(APIHandler{DB: m.db}),
		},
	}
}

func GetModule() lib.Module {
	m := &Module{}
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
