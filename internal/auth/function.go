package auth

import (
	"fmt"

	"awans.org/aft/internal/db"
	"github.com/google/uuid"
)

func unauthedCall(tx db.Tx, name string, args []interface{}) (result interface{}, err error) {
	function := tx.Ref(db.FunctionInterface.ID())

	funcRec, err := tx.Query(function,
		db.Filter(function, db.Eq("name", name)),
		db.Filter(function, db.Eq("funcType", uuid.UUID(db.RPC.ID()))),
	).OneRecord()
	if err == db.ErrNotFound {
		err = fmt.Errorf("%w: function %v not found", err, name)
		return
	} else if err != nil {
		return
	}
	f, err := tx.Schema().LoadFunction(funcRec)
	if err != nil {
		return
	}
	result, err = f.Call(tx.Context(), args)
	return
}

func AuthedCall(tx db.Tx, name string, args []interface{}) (result interface{}, err error) {
	_, isAuthedTx := tx.(*authedTx)
	if !isAuthedTx {
		_, isAuthedTx = tx.(*authedRWTx)
	}

	if !isAuthedTx {
		return unauthedCall(tx, name, args)
	}

	initialCtx := tx.Context()
	ctx := initialCtx
	role, ok := roleFromContext(ctx)
	if !ok {
		return nil, fmt.Errorf("No role in context calling %v\n", name)
	}
	function := tx.Ref(db.FunctionInterface.ID())
	functionRoles := tx.Ref(RoleModel.ID())
	currentRole := tx.Ref(RoleModel.ID())

	elevatedTx := Escalate(tx)
	qr, err := elevatedTx.Query(function,
		db.Filter(function, db.Eq("name", name)),
		db.Filter(function, db.Eq("funcType", uuid.UUID(db.RPC.ID()))),
		db.Filter(currentRole, db.EqID(role.ID())),
		db.Join(currentRole, function.Rel(ExecutableBy.Load(tx))),
		db.LeftJoin(functionRoles, function.Rel(FunctionRole.Load(tx))),
	).One()

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
		tx = tx.WithContext(ctx)
		tx, err = ActAsFunction(tx)
		if err != nil {
			return nil, err
		}
	}
	f, err := tx.Schema().LoadFunction(funcRec)
	if err != nil {
		return
	}
	result, err = f.Call(ctx, args)

	return
}
