package catalog

import (
	"awans.org/aft/internal/db"
	"awans.org/aft/internal/rpc"
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

func (m *Module) ProvideLiterals() []db.Literal {
	return []db.Literal{
		rpc.MakeRPC(
			db.MakeID("5dea9446-5b47-48a5-8bf8-7bd3678bab7c"),
			terminalRPC,
			nil,
		),
		rpc.MakeRPC(
			db.MakeID("e6082a9a-f181-46d7-8c22-e37d95b119dc"),
			lintRPC,
			nil,
		),
		rpc.MakeRPC(
			db.MakeID("e1106f29-48e2-4e3c-9f3e-9241c381a80c"),
			parseRPC,
			nil,
		),
	}
}

func GetModule() lib.Module {
	m := &Module{}
	return m
}
