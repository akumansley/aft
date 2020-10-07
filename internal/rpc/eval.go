package rpc

import (
	"context"

	"awans.org/aft/internal/db"
)

func eval(ctx context.Context, name string, args map[string]interface{}, tx db.RWTx) (interface{}, error) {
	function := tx.Ref(db.FunctionInterface.ID())

	// the munging for the enum EQ operation should be handled by matcher somehow
	// TODO add some kind of box that determines what functions are callable
	result, err := tx.Query(function, db.Filter(function, db.Eq("name", name))).OneRecord()
	if err != nil {
		return nil, err
	}

	f, err := tx.Schema().LoadFunction(result)
	if err != nil {
		return nil, err
	}
	return f.Call([]interface{}{ctx, args})
}
