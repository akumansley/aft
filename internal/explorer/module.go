package explorer

import (
	"io/ioutil"

	"awans.org/aft/internal/db"
	"awans.org/aft/internal/server/lib"
	"awans.org/aft/internal/starlark"
	"github.com/markbates/pkger"
)

type Module struct {
	lib.BlankModule
}

func init() {
	pkger.Include("/internal/explorer/reactForm.star")
	pkger.Include("/internal/explorer/validateForm.star")
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
	reactFormRPC := starlark.MakeStarlarkFunction(
		db.MakeID("d8179f1f-d94e-4b81-953b-6c170d3de9b7"),
		"reactForm",
		db.RPC,
		loadCode("/internal/explorer/reactForm.star"),
	)

	validateFormRPC := starlark.MakeStarlarkFunction(
		db.MakeID("d7633de5-9fa2-4409-a1b2-db96a59be52b"),
		"validateForm",
		db.RPC,
		loadCode("/internal/explorer/validateForm.star"),
	)

	return []db.FunctionL{
		reactFormRPC,
		validateFormRPC,
	}
}

func GetModule() lib.Module {
	m := &Module{}
	return m
}
