package rpcs

import (
	_ "embed"

	"awans.org/aft/internal/api/functions"
	"awans.org/aft/internal/auth"
	"awans.org/aft/internal/db"
	"awans.org/aft/internal/server/lib"
	"awans.org/aft/internal/starlark"
)

type Module struct {
	lib.BlankModule
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

	auth.Public.Functions = []db.FunctionL{
		loginRPC,
		meRPC,
	}
	auth.UserRoleL.Functions = []db.FunctionL{
		logoutRPC,
		meRPC,
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

//go:embed login.star
var loginStar string

var loginRPC = starlark.MakeStarlarkFunctionWithRole(
	db.MakeID("bf78428c-76ee-47d5-bc10-36788e1edede"),
	"login",
	1,
	db.RPC,
	loginStar,
	auth.LoginSystem,
)

//go:embed signup.star
var signupStar string

var signupRPC = starlark.MakeStarlarkFunction(
	db.MakeID("37371944-3728-42ba-8228-012d2f3702ad"),
	"signup",
	1,
	db.RPC,
	signupStar,
)

//go:embed me.star
var meStar string

var meRPC = starlark.MakeStarlarkFunction(
	db.MakeID("accda2d2-b217-4f05-bc2b-9a7f1ee8168a"),
	"me",
	1,
	db.RPC,
	meStar,
)

//go:embed logout.star
var logoutStar string

var logoutRPC = starlark.MakeStarlarkFunction(
	db.MakeID("86519f44-f1fe-4683-9647-91d0e02d5fd0"),
	"logout",
	1,
	db.RPC,
	logoutStar,
)

func (m *Module) ProvideFunctions() []db.FunctionL {
	return []db.FunctionL{
		auth.AuthenticateAs,
		auth.ClearAuthentication,
		loginRPC,
		signupRPC,
		meRPC,
		logoutRPC,
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
