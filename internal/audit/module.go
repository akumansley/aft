package audit

import (
	"awans.org/aft/internal/bus"
	"awans.org/aft/internal/oplog"
	"awans.org/aft/internal/server/lib"
)

type Module struct {
	lib.BlankModule
	b     *bus.EventBus
	audit oplog.OpLog
}

func (m *Module) ProvideRoutes() []lib.Route {
	return []lib.Route{
		lib.Route{
			Name:    "LogScan",
			Pattern: "log.scan",
			Handler: lib.ErrorHandler(LogScanHandler{bus: m.b, log: m.audit}),
		},
	}
}

func GetModule(b *bus.EventBus, log oplog.OpLog) lib.Module {
	return &Module{b: b, audit: log}
}
