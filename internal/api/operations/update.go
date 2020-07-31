package operations

import (
	"awans.org/aft/internal/db"
	"fmt"
)

func (op UpdateOperation) Apply(tx db.RWTx) (*db.QueryResult, error) {
	root := tx.Ref(op.ModelID)
	clauses := handleFindMany(tx, root, op.FindArgs)
	q := tx.Query(root, clauses...)
	outs := q.All()
	if len(outs) > 1 {
		return nil, fmt.Errorf("Found more than one record")
	}
	if len(outs) == 0 {
		return nil, nil
	}
	for _, no := range op.Nested {
		err := no.ApplyNested(tx, root, outs)
		if err != nil {
			return nil, err
		}
	}
	oldRec := outs[0]
	newRec, err := updateRecordFromData(oldRec.Record, op.Data)
	if err != nil {
		return nil, err
	}
	err = tx.Update(oldRec.Record, newRec)
	if err != nil {
		return nil, err
	}
	outs[0].Record = newRec
	err = tx.Commit()
	if err != nil {
		return nil, err
	}
	return outs[0], err
}

func (op NestedUpdateOperation) ApplyNested(tx db.RWTx, parent db.ModelRef, parents []*db.QueryResult) (err error) {
	outs, child := handleRelationshipWhere(tx, parent, parents, op.Relationship, op.Where)

	if len(outs) > 1 {
		return fmt.Errorf("Found more than one record")
	} else if len(outs) == 1 {
		for _, no := range op.Nested {
			err := no.ApplyNested(tx, child, outs)
			if err != nil {
				return err
			}
		}
		oldRec := outs[0].Record
		newRec, err := updateRecordFromData(oldRec, op.Data)
		err = tx.Update(oldRec, newRec)
		outs[0].Record = newRec
		if err != nil {
			return err
		}
	}
	return
}

func updateRecordFromData(oldRec db.Record, data map[string]interface{}) (db.Record, error) {
	newRec := oldRec.DeepCopy()
	for key, value := range data {
		err := newRec.Set(key, value)
		if err != nil {
			return nil, err
		}
		delete(data, key)
	}
	if len(data) != 0 {
		return nil, fmt.Errorf("Unused data in update")
	}
	return newRec, nil
}
