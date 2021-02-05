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

func (m *Module) ID() db.ID {
	return db.MakeID("3b092f9c-ff91-4543-9aa4-feaca9ff9c47")
}

func (m *Module) Name() string {
	return "auth"
}

func (m *Module) Package() string {
	return "awans.org/aft/internal/auth"
}

func GetModule(b *bus.EventBus) lib.Module {
	m := &Module{b: b}
	m.dbReadyHandler = func(event lib.DatabaseReady) {
		m.db = event.DB
	}
	return m
}

func (m *Module) ProvideInterfaces() []db.InterfaceL {
	return []db.InterfaceL{
		AuthKeyModel,
		PolicyModel,
		RoleModel,
		UserModel,
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
		FunctionRole,
		NativeFunctionRole,
		ExecutableBy,
		NativeFunctionExecutableBy,
		ModuleRoles,
		RoleModule,
	}
}

func (m *Module) ProvideFunctions() []db.FunctionL {
	return []db.FunctionL{
		CurrentUser,
		passwordValidator,
		emailAddressValidator,
		CheckPassword,
	}
}

func (m *Module) ProvideDatatypes() []db.DatatypeL {
	return []db.DatatypeL{
		Password,
		EmailAddress,
	}
}

func (m *Module) ProvideMiddleware() []lib.Middleware {
	return []lib.Middleware{
		makeAuthMiddleware(m.db),
	}
}
