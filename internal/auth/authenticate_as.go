package auth

import (
	"context"
	"time"

	"awans.org/aft/internal/db"
	"github.com/google/uuid"
)

var ttl = 30 * time.Minute

func authenticateAsFunc(args []interface{}) (result interface{}, err error) {
	ctx := args[0].(context.Context)
	id := args[1].(uuid.UUID)
	return
}

var AuthenticateAs = db.MakeNativeFunction(
	db.MakeID("e20ae44f-6a5e-4d25-ab13-de3bd7b7c392"),
	"authenticateAs",
	2,
	authenticateAsFunc,
)
