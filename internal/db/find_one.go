package db

import "awans.org/aft/internal/model"

func (op FindOneOperation) Apply(tx Tx) (st model.Record, err error) {
	return tx.FindOne(op.ModelName, op.UniqueQuery.Key, op.UniqueQuery.Val)
}
