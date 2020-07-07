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

func (m *Module) ProvideRecords() ([]db.Record, error) {
	datatypes := []db.Record{}
	r1 := db.RecordForModel(db.CodeModel)
	db.SaveCode(r1, EmailAddressValidator)
	datatypes = append(datatypes, r1)
	r2 := db.RecordForModel(db.CodeModel)
	db.SaveCode(r2, URLValidator)
	datatypes = append(datatypes, r2)
	r3 := db.RecordForModel(db.CodeModel)
	db.SaveCode(r3, PhoneValidator)
	datatypes = append(datatypes, r3)
	r4 := db.RecordForModel(db.DatatypeModel)
	db.SaveDatatype(r4, EmailAddress)
	datatypes = append(datatypes, r4)
	r5 := db.RecordForModel(db.DatatypeModel)
	db.SaveDatatype(r5, URL)
	datatypes = append(datatypes, r5)
	r6 := db.RecordForModel(db.DatatypeModel)
	db.SaveDatatype(r6, Phone)
	datatypes = append(datatypes, r6)
	return datatypes, nil
}
