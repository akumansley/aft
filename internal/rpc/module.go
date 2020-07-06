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
			Pattern: "/rpc/{name}",
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
	r, err := db.SaveRel(RPCCode)
	if err != nil {
		return records, err
	}
	records = append(records, r)

	rpcs := []db.Code{reactFormRPC, validateRPC, replRPC, parseRPC}
	ids := []uuid.UUID{
		uuid.MustParse("112197db-d9d6-46b7-9c9b-be4980562d95"),
		uuid.MustParse("865fbf7d-ce33-4e4c-bb7d-b4b5e1c82dca"),
		uuid.MustParse("f0626fc2-c6f9-4a93-be40-33a7fefaa548"),
		uuid.MustParse("34e086ea-b755-48e0-9a97-3073b2ca66ca")}
	for i, code := range rpcs {
		r1 := db.RecordForModel(db.CodeModel)
		db.SaveCode(r1, code)
		records = append(records, r1)

		r2 := db.RecordForModel(RPCModel)
		err = r2.Set("name", r1.MustGet("name"))
		if err != nil {
			return records, err
		}
		err = r2.Set("id", ids[i])
		if err != nil {
			return records, err
		}
		err = r2.SetFK("code", r1.ID())
		if err != nil {
			return records, err
		}
		records = append(records, r2)
	}
	return records, nil
}

func (m *Module) ProvideHandlers() []interface{} {
	return []interface{}{
		m.dbReadyHandler,
	}
}
