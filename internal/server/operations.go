package server

import (
	"awans.org/aft/internal/oplog"
	"awans.org/aft/internal/server/lib"
	"awans.org/aft/internal/server/operations"
)

func MakeOps(opLog oplog.OpLog) []lib.Operation {
	ops := []lib.Operation{
		lib.Operation{
			"LogScan",
			"log",
			"scan",
			operations.LogScanServer{Log: opLog},
			lib.None,
		},
		lib.Operation{
			"FindMany",
			"{object}",
			"findMany",
			operations.FindManyServer{},
			lib.Tx,
		},
		lib.Operation{
			"FindOne",
			"{object}",
			"findOne",
			operations.FindOneServer{},
			lib.Tx,
		},
		lib.Operation{
			"Create",
			"{object}",
			"create",
			operations.CreateServer{},
			lib.RWTx,
		},
	}
	return ops
}
