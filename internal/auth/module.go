package auth

import (
	"io/ioutil"

	"awans.org/aft/internal/bus"
	"awans.org/aft/internal/db"
	"awans.org/aft/internal/server/lib"
	"awans.org/aft/internal/starlark"
	"github.com/markbates/pkger"
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
	}
}

func init() {
	pkger.Include("/internal/auth/login.star")
}

func loadCode(path string) string {
	f, err := pkger.Open(path)
	if err != nil {
		panic(err)
	}
	b, err := ioutil.ReadAll(f)
	if err != nil {
		panic(err)
	}
	return string(b)
}

func (m *Module) ProvideFunctions() []db.FunctionL {
	loginRPC := starlark.MakeStarlarkFunction(
		db.MakeID("bf78428c-76ee-47d5-bc10-36788e1edede"),
		"login",
		2,
		loadCode("/internal/auth/login.star"),
	)
	return []db.FunctionL{
		loginRPC,
		AuthenticateAs,
	}
}
