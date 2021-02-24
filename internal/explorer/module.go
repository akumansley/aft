package explorer

import (
	_ "embed"

	"awans.org/aft/internal/db"
	"awans.org/aft/internal/server/lib"
	"awans.org/aft/internal/starlark"
)

type Module struct {
	lib.BlankModule
}

func (m *Module) ID() db.ID {
	return db.MakeID("4ec82123-e181-4121-ab4e-2dbfb0f3d09f")
}

func (m *Module) Name() string {
	return "explorer"
}

func (m *Module) Package() string {
	return "awans.org/aft/internal/explorer"
}

//go:embed reactForm.star
var reactFormStar string

var reactFormRPC = starlark.MakeStarlarkFunction(
	db.MakeID("d8179f1f-d94e-4b81-953b-6c170d3de9b7"),
	"reactForm",
	2,
	db.RPC,
	reactFormStar,
)

//go:embed validateForm.star
var validateFormStar string

var validateFormRPC = starlark.MakeStarlarkFunction(
	db.MakeID("d7633de5-9fa2-4409-a1b2-db96a59be52b"),
	"validateForm",
	2,
	db.RPC,
	validateFormStar,
)

func (m *Module) ProvideFunctions() []db.FunctionL {
	return []db.FunctionL{
		reactFormRPC,
		validateFormRPC,
	}
}

func GetModule() lib.Module {
	m := &Module{}
	return m
}
