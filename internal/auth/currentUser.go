package auth

import (
	"context"
	"errors"

	"awans.org/aft/internal/db"
)

func currentUserFunc(args []interface{}) (result interface{}, err error) {
	ctx := args[0].(context.Context)
	tx, ok := db.TxFromContext(ctx)

	user, ok := FromContext(tx, ctx)
	if ok {
		return user, nil
	}

	return nil, errors.New("Not signed in")
}

var CurrentUser = db.MakeNativeFunction(
	db.MakeID("5705529d-41ce-4b95-8e08-0725f4153d90"),
	"currentUser",
	1,
	currentUserFunc,
)
