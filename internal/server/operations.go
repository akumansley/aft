package server

import (
	"awans.org/aft/internal/oplog"
	"awans.org/aft/internal/server/db"
	"awans.org/aft/internal/server/operations"
	"net/http"
)

type Server interface {
	Parse(*http.Request) (interface{}, error)
	Serve(interface{}) (interface{}, error)
}

type Operation struct {
	Name    string
	Service string
	Method  string
	Server  Server
}

func MakeOps(db db.DB, opLog oplog.OpLog) []Operation {
	ops := []Operation{
		Operation{
			"LogScan",
			"log",
			"scan",
			operations.LogScanServer{Log: opLog},
		},
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
