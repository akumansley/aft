package catalog

import (
	"awans.org/aft/internal/db"
	"awans.org/aft/internal/server/lib"
)

type Module struct {
	lib.BlankModule
}

func (m *Module) ProvideFunctions() []db.FunctionL {

	return []db.FunctionL{
		terminalRPC,
		lintRPC,
		parseRPC,
	}
}

func GetModule() lib.Module {
	m := &Module{}
	return m
}
