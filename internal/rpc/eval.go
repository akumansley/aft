package rpc

import (
	"awans.org/aft/internal/db"
	"awans.org/aft/internal/starlark"
)

func eval(name string, args map[string]interface{}, tx db.RWTx) (interface{}, error) {
	r, err := tx.FindOne(RPCModel.ID, db.Eq("name", name))
	if err != nil {
		return nil, err
	}
	id, err := r.GetFK("code")
	if err != nil {
		return nil, err
	}
	c, err := tx.FindOne(db.CodeModel.ID, db.Eq("id", id))
	if err != nil {
		return nil, err
	}
	code, err := c.Get("code")
	if err != nil {
		return nil, err
	}
	sh := starlark.StarlarkFunctionHandle{Code: code.(string), Env: starlark.DBLib(tx)}
	return sh.Invoke(args)
}
