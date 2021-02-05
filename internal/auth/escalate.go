package auth

import (
	"context"

	"awans.org/aft/internal/db"
)

func withNoAuth(ctx context.Context) context.Context {
	return context.WithValue(ctx, noAuthKey, true)
}

func withAuth(ctx context.Context) context.Context {
	return context.WithValue(ctx, noAuthKey, false)
}

func shouldAuth(ctx context.Context) bool {
	noAuth, ok := ctx.Value(noAuthKey).(bool)
	if !ok {
		return true
	}
	return !noAuth
}

func Escalate(tx db.Tx) func() {
	oldCtx := tx.Context()
	newCtx := withNoAuth(oldCtx)
	tx.SetContext(newCtx)

	deescalate := func() {
		oldCtx2 := tx.Context()
		newCtx2 := withAuth(oldCtx2)
		tx.SetContext(newCtx2)
	}
	return deescalate
}
