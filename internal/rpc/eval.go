package rpc

import (
	"awans.org/aft/internal/db"
	"github.com/google/uuid"
)

func eval(name string, args map[string]interface{}, tx db.RWTx) (interface{}, error) {
	function := tx.Ref(db.FunctionInterface.ID())

	// the munging for the enum EQ operation should be handled by matcher somehow
	result, err := tx.Query(function, db.Filter(function, db.And(db.Eq("functionSignature", uuid.UUID(db.RPC.ID())), db.Eq("name", name)))).OneRecord()
	if err != nil {
		return nil, err
	}

	f, err := tx.Schema().LoadFunction(result)
	if err != nil {
		return nil, err
	}
	return f.Call(args)
}
