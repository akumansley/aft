package catalog

import (
	"errors"

	"awans.org/aft/internal/auth"
	"awans.org/aft/internal/db"
	"awans.org/aft/internal/server/lib"
	"github.com/google/uuid"
)

type Module struct {
	lib.BlankModule
}

func (m *Module) ID() db.ID {
	return db.MakeID("38e89e82-8721-41ea-860e-652684b58749")
}

func (m *Module) Name() string {
	return "catalog"
}

func (m *Module) Package() string {
	return "awans.org/aft/internal/catalog"
}

func (m *Module) ProvideFunctions() []db.FunctionL {
	return []db.FunctionL{
		terminalRPC,
		lintRPC,
		parseRPC,
	}
}

func (m *Module) ProvideHandlers() []interface{} {
	return []interface{}{
		initializeDefaultModule,
	}
}

func GetModule() lib.Module {
	m := &Module{}
	return m
}

var initializeDefaultModule = func(event lib.DatabaseReady) {
	appDB := event.DB
	tx := appDB.NewRWTx()
	auth.Escalate(tx)

	mods := tx.Ref(db.ModuleModel.ID())
	_, err := tx.Query(mods, db.Filter(mods, db.Eq("goPackage", ""))).OneRecord()
	if errors.Is(db.ErrNotFound, err) {
		newID := db.ID(uuid.New())
		modLit := db.MakeModule(newID, "app", "", []db.InterfaceL{}, []db.FunctionL{}, []db.DatatypeL{}, []db.ModuleLiteral{})
		appDB.AddLiteral(tx, modLit)
	}
	tx.Commit()
}
