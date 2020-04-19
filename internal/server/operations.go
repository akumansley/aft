package server

import (
	"awans.org/aft/internal/server/db"
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

func MakeOps(db db.DB) []Operation {
	ops := []Operation{
		Operation{
			"FindMany",
			"{object}",
			"findMany",
			operations.FindManyServer{DB: db},
		},
		Operation{
			"FindOne",
			"{object}",
			"findOne",
			operations.FindOneServer{DB: db},
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
