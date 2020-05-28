package api

import (
	"awans.org/aft/internal/bus"
	"awans.org/aft/internal/db"
	"awans.org/aft/internal/server/lib"
)

type Module struct {
	lib.BlankModule
	db db.DB
	b  *bus.EventBus
}

func (m Module) ProvideRoutes() []lib.Route {
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
			Handler: lib.ErrorHandler(CreateHandler{db: m.db, bus: m.b}),
		},
	}
}

func GetModule(db db.DB, b *bus.EventBus) lib.Module {
	return Module{db: db, b: b}
}
