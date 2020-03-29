package server

import (
	"fmt"
	"github.com/gorilla/mux"
)

func NewRouter() *mux.Router {
	router := mux.NewRouter()
	s := router.Methods("POST").Subrouter()
	for _, op := range operations {
		handler := Middleware(op)
		path := fmt.Sprintf("/api/%s.%s", op.Service, op.Method)
		s.Path(path).Name(op.Name).Handler(handler)
	}
	return router
}
