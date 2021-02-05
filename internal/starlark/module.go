package starlark

import (
	"awans.org/aft/internal/db"
	"awans.org/aft/internal/server/lib"
)

type Module struct {
	lib.BlankModule
}

func (m *Module) ID() db.ID {
	return db.MakeID("9b2e85f1-8914-453f-bd51-af5d45cc6f52")
}

func (m *Module) Name() string {
	return "starlark"
}

func (m *Module) Package() string {
	return "awans.org/aft/internal/starlark"
}

func (m *Module) ProvideFunctionLoaders() []db.FunctionLoader {
	sr := NewStarlarkRuntime()
	return []db.FunctionLoader{
		sr,
	}
}

func GetModule() lib.Module {
	m := &Module{}
	return m
}
