package auth

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
			Name:    "Login",
			Pattern: "/views/login",
			Handler: lib.ErrorHandler(LoginHandler{db: m.db, bus: m.b}),
		},
		lib.Route{
			Name:    "Signup",
			Pattern: "/views/signup",
			Handler: lib.ErrorHandler(SignupHandler{db: m.db, bus: m.b}),
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

func (m *Module) ProvideMiddleware() []lib.Middleware {
	return []lib.Middleware{
		makeAuthMiddleware(m.db),
	}
}

func (m *Module) ProvideModels() []db.ModelL {
	return []db.ModelL{
		AuthKeyModel,
		UserModel,
	}
}

func (m *Module) ProvideHandlers() []interface{} {
	return []interface{}{
		m.dbReadyHandler,
	}
}
