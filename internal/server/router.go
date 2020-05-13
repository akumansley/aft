package server

import (
	"awans.org/aft/internal/db"
	"awans.org/aft/internal/oplog"
	"awans.org/aft/internal/server/lib"
	"fmt"
	"github.com/gorilla/mux"
)

func NewRouter(ops []lib.Operation, db db.DB, log oplog.OpLog) *mux.Router {
	router := mux.NewRouter()
	s := router.Methods("POST").Subrouter()

	// defined in operations.go
	for _, op := range ops {
		handler := Middleware(op, db, log)
		path := fmt.Sprintf("/api/%s.%s", op.Service, op.Method)
		s.Path(path).Name(op.Name).Handler(handler)
	}
	return router
}
