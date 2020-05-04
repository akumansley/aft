package server

import (
	"awans.org/aft/internal/oplog"
	"awans.org/aft/internal/server/db"
	"awans.org/aft/internal/server/lib"
	"awans.org/aft/internal/server/operations"
	"context"
	"net/http"
)

func MakeOps(db db.DB, opLog oplog.OpLog) []lib.Operation {
	ops := []lib.Operation{
		lib.Operation{
			"LogScan",
			"log",
			"scan",
			operations.LogScanServer{Log: opLog},
			None,
		},
		lib.Operation{
			"FindMany",
			"{object}",
			"findMany",
			operations.FindManyServer{DB: db},
			Tx,
		},
		lib.Operation{
			"FindOne",
			"{object}",
			"findOne",
			operations.FindOneServer{DB: db},
			Tx,
		},
		lib.Operation{
			"Create",
			"{object}",
			"create",
			operations.CreateServer{DB: db},
			RWTx,
		},
	}
	return ops
}
