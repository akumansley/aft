package api

import (
	"awans.org/aft/internal/db"
)

type FindOneOperation struct {
	ModelID     db.ModelID
	UniqueQuery UniqueQuery
}

func (op FindOneOperation) Apply(tx db.Tx) (st db.Record, err error) {
	// TODO handle FK?
	return tx.FindOne(op.ModelID, db.Eq(op.UniqueQuery.Key, op.UniqueQuery.Val))
}
