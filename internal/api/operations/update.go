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
	for _, no := range op.Nested {
		err := no.ApplyNested(tx, root, root, outs, clauses)
		if err != nil {
			return nil, err
		}
	}
	err = tx.Commit()
	if err != nil {
		return nil, err
	}
	return outs[0], err
}

func (op NestedUpdateOperation) ApplyNested(tx db.RWTx, root db.ModelRef, parent db.ModelRef, parents []*db.QueryResult, clauses []db.QueryClause) (err error) {
	cls, child := handleRelationshipWhere(tx, parent, op.Relationship, op.Where)
	clauses = append(clauses, cls...)
	q := tx.Query(root, clauses...)
	outs := getEdgeResults(parents, q.All())

	if len(outs) > 1 {
		return fmt.Errorf("Found more than one record")
	} else if len(outs) == 1 {
		oldRec := outs[0].Record
		newRec, err := updateRecordFromData(oldRec, op.Data)
		err = tx.Update(oldRec, newRec)
		outs[0].Record = newRec
		if err != nil {
			return err
		}
		for _, no := range op.Nested {
			err := no.ApplyNested(tx, root, child, outs, clauses)
			if err != nil {
				return err
			}
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
