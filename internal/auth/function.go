package auth

import (
	"fmt"

	"awans.org/aft/internal/db"
	"github.com/google/uuid"
)

func AuthedCall(tx db.Tx, name string, args []interface{}) (result interface{}, err error) {
	initialCtx := tx.Context()
	ctx := initialCtx
	role, ok := roleFromContext(ctx)
	if !ok {
		return nil, fmt.Errorf("No role in context calling %v\n", name)
	}
	function := tx.Ref(db.FunctionInterface.ID())
	functionRoles := tx.Ref(RoleModel.ID())
	currentRole := tx.Ref(RoleModel.ID())

	deescalate := Escalate(tx)
	qr, err := tx.Query(function,
		db.Filter(function, db.Eq("name", name)),
		db.Filter(function, db.Eq("funcType", uuid.UUID(db.RPC.ID()))),
		db.Filter(currentRole, db.EqID(role.ID())),
		db.Join(currentRole, function.Rel(ExecutableBy.Load(tx))),
		db.LeftJoin(functionRoles, function.Rel(FunctionRole.Load(tx))),
	).One()
	deescalate()

	if err == db.ErrNotFound {
		err = fmt.Errorf("%w: function %v not found", err, name)
		return
	} else if err != nil {
		return
	}

	funcRec := qr.Record
	roleQR := qr.GetChildRelOne(FunctionRole.Load(tx))
	if roleQR != nil {
		role := roleQR.Record
		ctx = withFunctionRole(ctx, role)
		tx.SetContext(ctx)
	}
	ActAsFunction(tx)
	f, err := tx.Schema().LoadFunction(funcRec)
	if err != nil {
		return
	}
	result, err = f.Call(ctx, args)

	// no longer running as this func
	if roleQR != nil {
		tx.SetContext(initialCtx)
	}
	return
}
