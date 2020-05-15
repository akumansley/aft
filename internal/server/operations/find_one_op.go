package operations

import (
	"awans.org/aft/internal/db"
	"awans.org/aft/internal/model"
)

type FindOneOperation struct {
	ModelName   string
	UniqueQuery UniqueQuery
}

func (op FindOneOperation) Apply(tx db.Tx) (st model.Record, err error) {
	return tx.FindOne(op.ModelName, op.UniqueQuery.Key, op.UniqueQuery.Val)
}
