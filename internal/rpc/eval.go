package rpc

import (
	"awans.org/aft/internal/db"
	"fmt"
	"github.com/google/uuid"
)

func eval(name string, args map[string]interface{}, tx db.RWTx) (interface{}, error) {
	function := tx.Ref(db.FunctionInterface.ID())

	// the munging for the enum EQ operation should be handled by matcher somehow
	results := tx.Query(function, db.Filter(function, db.And(db.Eq("functionSignature", uuid.UUID(db.RPC.ID())), db.Eq("name", name)))).All()
	if len(results) > 1 {
		err := fmt.Errorf("multiple rpcs named %v: %v", name, results)
		panic(err)
	}
	if len(results) == 0 {
		err := fmt.Errorf("no rpcs named %v", name)
		panic(err)
	}
	result := results[0]
	f, err := tx.Schema().LoadFunction(result.Record)
	if err != nil {
		return nil, err
	}
	return f.Call(args)
}
