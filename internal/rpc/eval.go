package rpc

import (
	"awans.org/aft/internal/db"
	"awans.org/aft/internal/starlark"
)

func eval(name string, args map[string]interface{}, tx db.RWTx) (interface{}, error) {
	rpcs := tx.Ref(RPCModel.ID)
	codes := tx.Ref(db.CodeModel.ID)
	results := tx.Query(rpcs).Join(codes, rpcs.Rel(RPCCode)).Filter(rpcs, db.Eq("name", name)).All()
	if len(results) != 1 {
		panic("multiple RPCs with the same name")
	}
	result := results[0]

	code := result.ToOne["code"].Record
	codeString, err := code.Get("code")
	if err != nil {
		return nil, err
	}
	sh := starlark.StarlarkFunctionHandle{Code: codeString.(string), Env: starlark.DBLib(tx)}
	return sh.Invoke(args)
}
