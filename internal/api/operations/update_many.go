package operations

import (
	"awans.org/aft/internal/db"
)

func (op UpdateManyOperation) Apply(tx db.RWTx) (int, error) {
	for i, _ := range op.Old {
		err := tx.Update(op.Old[i], op.New[i])
		if err != nil {
			return 0, err
		}
	}
	err := tx.Commit()
	if err != nil {
		return 0, err
	}
	return len(op.Old), nil
}

func (op NestedUpdateManyOperation) ApplyNested(tx db.RWTx, parent db.Record) (err error) {
	for i, _ := range op.Old {
		err := tx.Update(op.Old[i], op.New[i])
		if err != nil {
			return err
		}
	}
	err = tx.Commit()
	if err != nil {
		return err
	}
	return nil
}
