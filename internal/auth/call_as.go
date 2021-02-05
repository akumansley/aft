package auth

import (
	"fmt"

	"awans.org/aft/internal/db"
)

func ActAsFunction(tx db.Tx) error {
	oldCtx := tx.Context()
	functionRole, ok := functionRoleFromContext(oldCtx)
	if !ok {
		return fmt.Errorf("No function role in context")
	}
	newCtx := withRole(oldCtx, functionRole)
	tx.SetContext(newCtx)

	return nil
}

func ActAsUser(tx db.Tx) error {
	oldCtx := tx.Context()
	userRole, ok := userRoleFromContext(oldCtx)
	if !ok {
		return fmt.Errorf("No user role in context")
	}
	newCtx := withRole(oldCtx, userRole)
	tx.SetContext(newCtx)
	return nil
}
