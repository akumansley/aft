package server

import (
	"awans.org/aft/internal/server/operations"
	"net/http"
)

type Server interface {
	Parse(*http.Request) interface{}
	Serve(http.ResponseWriter, interface{})
}

type Operation struct {
	Name    string
	Service string
	Method  string
	Server  Server
}

func MakeOps(db DB) {
	ops := []Operation{
		Operation{
			"Query",
			"{object}",
			"query",
			operations.QueryServer{DB: db},
		},
		Operation{
			"Create",
			"{object}",
			"create",
			operations.CreateServer{DB: db},
		},
	}
	return ops
}
