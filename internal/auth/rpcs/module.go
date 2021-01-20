package rpcs

import (
	"io/ioutil"

	"awans.org/aft/internal/auth"
	"awans.org/aft/internal/db"
	"awans.org/aft/internal/rpc"
	"awans.org/aft/internal/server/lib"
	"awans.org/aft/internal/starlark"
	"github.com/markbates/pkger"
)

type Module struct {
	lib.BlankModule
	loginRPC, signupRPC, meRPC db.FunctionL
}

func GetModule() lib.Module {
	m := &Module{}
	m.loginRPC = starlark.MakeStarlarkFunction(
		db.MakeID("bf78428c-76ee-47d5-bc10-36788e1edede"),
		"login",
		2,
		loadCode("/internal/auth/rpcs/login.star"),
	)
	m.signupRPC = starlark.MakeStarlarkFunction(
		db.MakeID("37371944-3728-42ba-8228-012d2f3702ad"),
		"signup",
		2,
		loadCode("/internal/auth/rpcs/signup.star"),
	)
	m.meRPC = starlark.MakeStarlarkFunction(
		db.MakeID("accda2d2-b217-4f05-bc2b-9a7f1ee8168a"),
		"me",
		2,
		loadCode("/internal/auth/rpcs/me.star"),
	)

	return m
}

func (m *Module) ProvideLiterals() []db.Literal {
	return []db.Literal{
		rpc.MakeRPC(
			db.MakeID("f2edb14e-19c9-493a-93da-bab2ee865907"),
			m.loginRPC,
			&auth.System,
		),
		rpc.MakeRPC(
			db.MakeID("7451ec9b-9f37-4681-b88e-e1762319ae80"),
			m.signupRPC,
			nil,
		),
		rpc.MakeRPC(
			db.MakeID("82b374b7-7453-4934-8ddb-59e6aa77ecd9"),
			m.meRPC,
			nil,
		),
	}
}

func init() {
	pkger.Include("/internal/auth/rpcs/login.star")
	pkger.Include("/internal/auth/rpcs/signup.star")
	pkger.Include("/internal/auth/rpcs/me.star")
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
	}
}
