package auth

import (
	"io/ioutil"

	"awans.org/aft/internal/bus"
	"awans.org/aft/internal/db"
	"awans.org/aft/internal/rpc"
	"awans.org/aft/internal/server/lib"
	"awans.org/aft/internal/starlark"
	"github.com/markbates/pkger"
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
	m.loginRPC = starlark.MakeStarlarkFunction(
		db.MakeID("bf78428c-76ee-47d5-bc10-36788e1edede"),
		"login",
		2,
		loadCode("/internal/auth/login.star"),
	)
	m.signupRPC = starlark.MakeStarlarkFunction(
		db.MakeID("37371944-3728-42ba-8228-012d2f3702ad"),
		"signup",
		2,
		loadCode("/internal/auth/signup.star"),
	)
	m.meRPC = starlark.MakeStarlarkFunction(
		db.MakeID("accda2d2-b217-4f05-bc2b-9a7f1ee8168a"),
		"me",
		2,
		loadCode("/internal/auth/me.star"),
	)

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
		rpc.MakeRPC(
			db.MakeID("f2edb14e-19c9-493a-93da-bab2ee865907"),
			m.loginRPC,
		),
		rpc.MakeRPC(
			db.MakeID("7451ec9b-9f37-4681-b88e-e1762319ae80"),
			m.signupRPC,
		),
		rpc.MakeRPC(
			db.MakeID("82b374b7-7453-4934-8ddb-59e6aa77ecd9"),
			m.meRPC,
		),
	}
}

func init() {
	pkger.Include("/internal/auth/login.star")
	pkger.Include("/internal/auth/signup.star")
	pkger.Include("/internal/auth/me.star")
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
	return []db.FunctionL{
		m.loginRPC,
		m.signupRPC,
		m.meRPC,
		AuthenticateAs,
		CurrentUser,
		passwordValidator,
		HashPassword,
	}
}

func (m *Module) ProvideDatatypes() []db.DatatypeL {
	return []db.DatatypeL{
		Password,
	}
}
