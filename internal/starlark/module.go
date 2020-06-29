package starlark

import (
	"awans.org/aft/internal/db"
	"awans.org/aft/internal/server/lib"
)

type Module struct {
	lib.BlankModule
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
