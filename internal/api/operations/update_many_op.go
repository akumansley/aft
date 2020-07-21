package operations

import (
	"awans.org/aft/internal/db"
)

type UpdateManyOperation struct {
	Old []db.Record
	New []db.Record
}

func (op UpdateManyOperation) Apply(tx db.RWTx) (int, error) {
	for i, _ := range op.Old {
		err := tx.Update(op.Old[i], op.New[i])
		if err != nil {
			return 0, nil
		}
	}
	err := tx.Commit()
	if err != nil {
		return 0, err
	}
	return len(op.Old), nil
}
