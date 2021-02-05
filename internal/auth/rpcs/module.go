package rpcs

import (
	"io/ioutil"

	"awans.org/aft/internal/api/functions"
	"awans.org/aft/internal/auth"
	"awans.org/aft/internal/db"
	"awans.org/aft/internal/server/lib"
	"awans.org/aft/internal/starlark"
	"github.com/markbates/pkger"
)

type Module struct {
	lib.BlankModule
	loginRPC, signupRPC, meRPC, logoutRPC db.FunctionL
}

func (m *Module) ID() db.ID {
	return db.MakeID("b693ff27-5409-42f8-bf16-2c3e86f886a9")
}

func (m *Module) Name() string {
	return "loginSystem"
}

func (m *Module) Package() string {
	return "awans.org/aft/internal/auth/rpcs"
}

func GetModule() lib.Module {
	m := &Module{}
	m.loginRPC = starlark.MakeStarlarkFunctionWithRole(
		db.MakeID("bf78428c-76ee-47d5-bc10-36788e1edede"),
		"login",
		1,
		db.RPC,
		loadCode("/internal/auth/rpcs/login.star"),
		auth.LoginSystem,
	)
	m.signupRPC = starlark.MakeStarlarkFunction(
		db.MakeID("37371944-3728-42ba-8228-012d2f3702ad"),
		"signup",
		1,
		db.RPC,
		loadCode("/internal/auth/rpcs/signup.star"),
	)
	m.meRPC = starlark.MakeStarlarkFunction(
		db.MakeID("accda2d2-b217-4f05-bc2b-9a7f1ee8168a"),
		"me",
		1,
		db.RPC,
		loadCode("/internal/auth/rpcs/me.star"),
	)
	m.logoutRPC = starlark.MakeStarlarkFunction(
		db.MakeID("86519f44-f1fe-4683-9647-91d0e02d5fd0"),
		"logout",
		1,
		db.RPC,
		loadCode("/internal/auth/rpcs/logout.star"),
	)
	auth.Public.Functions = []db.FunctionL{
		m.loginRPC,
		m.meRPC,
	}
	auth.UserRoleL.Functions = []db.FunctionL{
		m.logoutRPC,
		m.meRPC,
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

	return m
}

func init() {
	pkger.Include("/internal/auth/rpcs/login.star")
	pkger.Include("/internal/auth/rpcs/signup.star")
	pkger.Include("/internal/auth/rpcs/me.star")
	pkger.Include("/internal/auth/rpcs/logout.star")
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
		auth.AuthenticateAs,
		auth.ClearAuthentication,
		m.loginRPC,
		m.signupRPC,
		m.meRPC,
		m.logoutRPC,
	}
}

func (m *Module) ProvideLiterals() []db.Literal {
	return []db.Literal{
		auth.Public,
		auth.UserUserPolicy,
		auth.UserRoleL,
		auth.LoginUserPolicy,
		auth.LoginAuthKeyPolicy,
		auth.LoginSystem,
	}
}
