package auth

import (
	"fmt"

	"awans.org/aft/internal/db"
)

func ActAsFunction(tx db.Tx) (db.Tx, error) {
	oldCtx := tx.Context()
	functionRole, ok := functionRoleFromContext(oldCtx)
	if !ok {
		return nil, fmt.Errorf("No function role in context")
	}
	newCtx := withRole(oldCtx, functionRole)
	return tx.WithContext(newCtx), nil
}

func ActAsUser(tx db.Tx) (db.Tx, error) {
	oldCtx := tx.Context()
	userRole, ok := userRoleFromContext(oldCtx)
	if !ok {
		return nil, fmt.Errorf("No user role in context")
	}
	newCtx := withRole(oldCtx, userRole)
	return tx.WithContext(newCtx), nil
}
