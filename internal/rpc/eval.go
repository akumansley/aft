package rpc

import (
	"awans.org/aft/internal/db"
)

func eval(name string, args map[string]interface{}, tx db.RWTx) (interface{}, error) {
	rpcs := tx.Ref(RPCModel.ID())
	function := tx.Ref(db.FunctionInterface.ID())
	results := tx.Query(rpcs).Join(function, rpcs.Rel(RPCCode)).Filter(rpcs, db.Eq("name", name)).All()
	if len(results) != 1 {
		panic("multiple RPCs with the same name")
	}
	result := results[0]
	f, err := tx.Schema().LoadFunction(result.Record)
	if err != nil {
		return nil, err
	}
	return f.Call(args)
}
