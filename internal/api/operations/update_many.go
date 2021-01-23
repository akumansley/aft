package operations

import (
	"awans.org/aft/internal/db"
)

func (op UpdateManyOperation) Apply(tx db.RWTx) (int, error) {
	root := tx.Ref(op.ModelID)
	clauses := HandleWhere(tx, root, op.Where)
	q := tx.Query(root, clauses...)
	oldRecs := q.All()

	for _, oldRec := range oldRecs {
		newRec, err := updateRecordFromData(tx, oldRec.Record, op.Data)
		if err != nil {
			return 0, err
		}
		err = tx.Update(oldRec.Record, newRec)
		oldRec.Record = newRec
		if err != nil {
			return 0, err
		}
	}
	for _, no := range op.Nested {
		err := no.ApplyNested(tx, root, oldRecs)
		if err != nil {
			return 0, err
		}
	}

	return len(oldRecs), nil
}

func (op NestedUpdateManyOperation) ApplyNested(tx db.RWTx, parent db.ModelRef, parents []*db.QueryResult) (err error) {
	outs, child := handleRelationshipWhere(tx, parent, parents, op.Relationship, op.Where)

	for _, no := range op.Nested {
		err := no.ApplyNested(tx, child, outs)
		if err != nil {
			return err
		}
	}
	for i, oldRec := range outs {
		newRec, err := updateRecordFromData(tx, oldRec.Record, op.Data)
		if err != nil {
			return err
		}
		err = tx.Update(oldRec.Record, newRec)
		outs[i].Record = newRec
		if err != nil {
			return err
		}
	}
	return
}
