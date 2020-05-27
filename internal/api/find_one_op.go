package api

import (
	"awans.org/aft/internal/db"
)

type FindOneOperation struct {
	ModelName   string
	UniqueQuery UniqueQuery
}

func (op FindOneOperation) Apply(tx db.Tx) (st db.Record, err error) {
	// TODO handle FK?
	return tx.FindOne(op.ModelName, db.Eq(op.UniqueQuery.Key, op.UniqueQuery.Val))
}
