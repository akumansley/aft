package operations

import (
	"awans.org/aft/internal/db"
)

type UpdateOperation struct {
	Old db.Record
	New db.Record
}

func (op UpdateOperation) Apply(tx db.RWTx) (db.Record, error) {
	err := tx.Update(op.Old, op.New)
	if err != nil {
		return nil, err
	}
	err = tx.Commit()
	if err != nil {
		return nil, err
	}
	return op.New, nil
}
