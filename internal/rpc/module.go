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

func (m *Module) ProvideModels() []db.Model {
	return []db.Model{
		RPCModel,
	}
}

func (m *Module) ProvideRecords() ([]db.Record, error) {
	records := []db.Record{}
	r1, err := db.SaveRel(RPCCode)
	if err != nil {
		return records, err
	}
	records = append(records, r1)
	r2 := db.RecordForModel(db.CodeModel)
	db.SaveCode(r2, reactFormRPC)
	records = append(records, r2)
	r3 := db.RecordForModel(RPCModel)
	err = r3.Set("name", "reactForm")
	if err != nil {
		return records, err
	}
	err = r3.Set("id", uuid.MustParse("112197db-d9d6-46b7-9c9b-be4980562d95"))
	if err != nil {
		return records, err
	}
	err = r3.SetFK("code", r2.ID())
	if err != nil {
		return records, err
	}
	records = append(records, r3)
	return records, nil
}

func (m *Module) ProvideHandlers() []interface{} {
	return []interface{}{
		m.dbReadyHandler,
	}
}
