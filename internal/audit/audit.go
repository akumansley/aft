package audit

import (
	"awans.org/aft/internal/bus"
	"awans.org/aft/internal/oplog"
	"awans.org/aft/internal/server/lib"
	"github.com/json-iterator/go"
	"net/http"
)

type Module struct {
	lib.BlankModule
	b     *bus.EventBus
	audit oplog.OpLog
}

func (m Module) ProvideRoutes() []lib.Route {
	return []lib.Route{
		lib.Route{
			Name:    "LogScan",
			Pattern: "log.scan",
			Handler: errorHandler(LogScanHandler{bus: m.b, log: m.audit}),
		},
	}
}

func GetModule(b *bus.EventBus, log oplog.OpLog) lib.Module {
	return Module{b: b, audit: log}
}

type apiHandler interface {
	ServeHTTP(http.ResponseWriter, *http.Request) error
}

type ErrorResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func errorHandler(inner apiHandler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := inner.ServeHTTP(w, r)
		if err != nil {
			er := ErrorResponse{
				Code:    "serve-error",
				Message: err.Error(),
			}
			bytes, _ := jsoniter.Marshal(&er)
			status := http.StatusBadRequest

			_, _ = w.Write(bytes)
			w.WriteHeader(status)
		}
	})
}
