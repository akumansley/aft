package api

import (
	"awans.org/aft/internal/db"
	"github.com/google/uuid"
)

type FindOneOperation struct {
	ModelID     uuid.UUID
	UniqueQuery UniqueQuery
}

func (op FindOneOperation) Apply(tx db.Tx) (st db.Record, err error) {
	// TODO handle FK?
	return tx.FindOne(op.ModelID, db.Eq(op.UniqueQuery.Key, op.UniqueQuery.Val))
}
