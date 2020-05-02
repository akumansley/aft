package server

import (
	"awans.org/aft/internal/oplog"
	"fmt"
	"github.com/gorilla/mux"
)

func NewRouter(ops []Operation, log oplog.OpLog) *mux.Router {
	router := mux.NewRouter()
	s := router.Methods("POST").Subrouter()

	// defined in operations.go
	for _, op := range ops {
		handler := Middleware(op, log)
		path := fmt.Sprintf("/api/%s.%s", op.Service, op.Method)
		s.Path(path).Name(op.Name).Handler(handler)
	}
	return router
}
