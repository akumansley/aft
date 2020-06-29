package bizdatatypes

import (
	"awans.org/aft/internal/bus"
	"awans.org/aft/internal/db"
	"awans.org/aft/internal/server/lib"
)

type Module struct {
	lib.BlankModule
	bus            *bus.EventBus
	db             db.DB
	dbReadyHandler interface{}
}

func GetModule(b *bus.EventBus) lib.Module {
	m := &Module{bus: b}
	m.dbReadyHandler = func(event lib.DatabaseReady) {
		m.db = event.Db
	}
	return m
}

func (m *Module) ProvideFunctions() []db.FunctionL {
	return []db.FunctionL{
		EmailAddressValidator,
		URLValidator,
	}
}

func (m *Module) ProvideDatatypes() []db.DatatypeL {
	return []db.DatatypeL{
		EmailAddress,
		URL,
	}
}
