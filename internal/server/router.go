package server

import (
	"fmt"
	"github.com/gorilla/mux"
)

func NewRouter(ops []Operation) *mux.Router {
	router := mux.NewRouter()
	s := router.Methods("POST").Subrouter()

	// defined in operations.go
	for _, op := range ops {
		handler := Middleware(op)
		path := fmt.Sprintf("/api/%s.%s", op.Service, op.Method)
		s.Path(path).Name(op.Name).Handler(handler)
	}
	return router
}
