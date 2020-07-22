package operations

import (
	"awans.org/aft/internal/db"
)

func (op UpdateManyOperation) Apply(tx db.RWTx) (int, error) {
	fm := FindManyOperation{ModelID: op.ModelID, FindManyArgs: FindManyArgs{Where: op.Where}}
	oldRecs, err := fm.Apply(tx)
	if err != nil {
		return 0, err
	}
	for _, oldRec := range oldRecs {
		newRec, err := updateRecordFromData(oldRec.Record, op.Data)
		if err != nil {
			return 0, err
		}
		err = tx.Update(oldRec.Record, newRec)
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
