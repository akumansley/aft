package audit

import (
	"awans.org/aft/internal/bus"
	"awans.org/aft/internal/db"
	"awans.org/aft/internal/oplog"
	"awans.org/aft/internal/server/lib"
)

type Module struct {
	lib.BlankModule
	b          *bus.EventBus
	requestLog oplog.OpLog
	dbLog      oplog.OpLog
}

func (m *Module) ProvideFunctions() []db.FunctionL {
	return []db.FunctionL{
		makeScanFunction(map[string]oplog.OpLog{
			"request": m.requestLog,
			"db":      m.dbLog,
		}),
	}
}

func GetModule(b *bus.EventBus, dbLog, requestLog oplog.OpLog) lib.Module {
	return &Module{b: b, dbLog: dbLog, requestLog: requestLog}
}

func (m *Module) ProvideHandlers() []interface{} {
	return []interface{}{
		makeSaveRequestsHandler(m.requestLog),
	}
}

func makeSaveRequestsHandler(log oplog.OpLog) interface{} {
	handler := func(event lib.ParseRequest) {
		log.Log(event.Request)
	}
	return handler
}
