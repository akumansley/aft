package operations

import (
	"awans.org/aft/internal/db"
)

type FindOneOperation struct {
	ModelID     db.ID
	UniqueQuery UniqueQuery
}

func (op FindOneOperation) Apply(tx db.Tx) (st db.Record, err error) {
	t := tx.Ref(op.ModelID)
	return tx.Query(t, db.Filter(t, db.Eq(op.UniqueQuery.Key, op.UniqueQuery.Val))).OneRecord()
}
