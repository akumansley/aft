package audit

import (
	"net/http"

	"awans.org/aft/internal/bus"
	"awans.org/aft/internal/db"
	"awans.org/aft/internal/oplog"
	"awans.org/aft/internal/rpc"
	"awans.org/aft/internal/server/lib"
)

type Module struct {
	lib.BlankModule
	b          *bus.EventBus
	requestLog oplog.OpLog
	scanFunc   db.FunctionL
}

func (m *Module) ProvideFunctions() []db.FunctionL {
	return []db.FunctionL{m.scanFunc}
}

func (m *Module) ProvideLiterals() []db.Literal {
	return []db.Literal{
		rpc.MakeRPC(
			db.MakeID("9e63486a-a010-4ff8-b2ad-0a33eba43fc2"),
			m.scanFunc,
			nil,
		),
	}
}

func GetModule(b *bus.EventBus, dbLog oplog.OpLog) lib.Module {
	requestLog := oplog.GobLog(oplog.NewMemLog())
	scanFunc := makeScanFunction(map[string]oplog.OpLog{
		"request": requestLog,
		"db":      dbLog,
	})
	return &Module{b: b, scanFunc: scanFunc, requestLog: requestLog}
}

func (m *Module) ProvideMiddleware() []lib.Middleware {
	return []lib.Middleware{
		makeAuditMiddleware(m.b),
	}
}

func makeAuditMiddleware(b *bus.EventBus) lib.Middleware {
	return func(inner http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := bus.WithBus(r.Context(), b)
			r = r.Clone(ctx)
			inner.ServeHTTP(w, r)
		})
	}
}

func (m *Module) ProvideHandlers() []interface{} {
	return []interface{}{
		makeSaveRequestsHandler(m.requestLog),
	}
}

func makeSaveRequestsHandler(log oplog.OpLog) interface{} {
	handler := func(event lib.ParseRequest) {
		err := log.Log(event.Request)
		if err != nil {
			panic(err)
		}
	}
	return handler
}
