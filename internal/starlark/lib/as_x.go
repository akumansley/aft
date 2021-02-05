package lib

import (
	"context"
	"errors"

	"awans.org/aft/internal/auth"
	"awans.org/aft/internal/db"
	"go.starlark.net/starlark"
)

var asSelf = starlark.NewBuiltin("asSelf", setRoleSelf)
var asUser = starlark.NewBuiltin("asUser", setRoleUser)

func setRoleSelf(thread *starlark.Thread, b *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	ctx := thread.Local("ctx").(context.Context)
	tx, ok := db.TxFromContext(ctx)
	if !ok {
		return nil, errors.New("No tx in context")
	}
	auth.ActAsFunction(tx)

	return starlark.None, nil
}

func setRoleUser(thread *starlark.Thread, b *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	ctx := thread.Local("ctx").(context.Context)
	tx, ok := db.TxFromContext(ctx)
	if !ok {
		return nil, errors.New("No tx in context")
	}
	auth.ActAsUser(tx)

	return starlark.None, nil
}
