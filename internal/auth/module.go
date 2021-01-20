package auth

import (
	"awans.org/aft/internal/bus"
	"awans.org/aft/internal/db"
	"awans.org/aft/internal/server/lib"
)

type Module struct {
	lib.BlankModule
	db                         db.DB
	b                          *bus.EventBus
	dbReadyHandler             interface{}
	loginRPC, signupRPC, meRPC db.FunctionL
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
		RoleModel,
		PolicyModel,
	}
}

func (m *Module) ProvideHandlers() []interface{} {
	return []interface{}{
		m.dbReadyHandler,
		initializeAuthKey,
	}
}

func (m *Module) ProvideLiterals() []db.Literal {
	return []db.Literal{
		Public,
		System,
	}
}

func (m *Module) ProvideFunctions() []db.FunctionL {
	return []db.FunctionL{
		AuthenticateAs,
		CurrentUser,
		passwordValidator,
		CheckPassword,
	}
}

func (m *Module) ProvideDatatypes() []db.DatatypeL {
	return []db.DatatypeL{
		Password,
	}
}
