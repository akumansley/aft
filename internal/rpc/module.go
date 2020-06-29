package rpc

import (
	"awans.org/aft/internal/bus"
	"awans.org/aft/internal/db"
	"awans.org/aft/internal/server/lib"
	"github.com/google/uuid"
)

type Module struct {
	lib.BlankModule
	db             db.DB
	bus            *bus.EventBus
	dbReadyHandler interface{}
}

func (m Module) ProvideRoutes() []lib.Route {
	return []lib.Route{
		lib.Route{
			Name:    "RPC",
			Pattern: "/views/rpc/{name}",
			Handler: lib.ErrorHandler(RPCHandler{db: m.db, bus: m.bus}),
		},
	}
}

func GetModule(b *bus.EventBus) lib.Module {
	m := &Module{bus: b}
	m.dbReadyHandler = func(event lib.DatabaseReady) {
		m.db = event.Db
	}
	return m
}

func (m *Module) ProvideModels() []db.ModelL {
	return []db.ModelL{
		RPCModel,
	}
}

func (m *Module) ProvideRelationships() []db.Relationship {
	return []db.Relationship{
		RPCCode,
	}
}

func (m *Module) ProvideFunctions() []db.FunctionL {
	return []db.FunctionL{
		reactFormRPC,
	}
}

func (m *Module) ProvideRecords(tx db.RWTx) (err error) {
	r3 := db.RecordForModel(RPCModel)
	err = r3.Set("name", "reactForm")
	if err != nil {
		return
	}
	err = r3.Set("id", uuid.MustParse("112197db-d9d6-46b7-9c9b-be4980562d95"))
	if err != nil {
		return
	}
	tx.Connect(r3.ID(), reactFormRPC.ID(), RPCCode.ID())
	return nil
}

func (m *Module) ProvideHandlers() []interface{} {
	return []interface{}{
		m.dbReadyHandler,
	}
}
