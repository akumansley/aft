package auth

import (
	"context"
	"errors"

	"awans.org/aft/internal/db"
)

func currentUserFunc(ctx context.Context, args []interface{}) (result interface{}, err error) {
	user, ok := UserFromContext(ctx)
	if ok {
		return user, nil
	}

	return nil, errors.New("Not signed in")
}

var CurrentUser = db.MakeNativeFunction(
	db.MakeID("5705529d-41ce-4b95-8e08-0725f4153d90"),
	"currentUser",
	0,
	db.Internal,
	currentUserFunc,
)
