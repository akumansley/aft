package operations

import (
	"awans.org/aft/internal/db"
)

func (op UpdateManyOperation) Apply(tx db.RWTx) (int, error) {
	root := tx.Ref(op.ModelID)
	clauses := handleWhere(tx, root, op.Where)
	q := tx.Query(root, clauses...)
	oldRecs := q.All()

	for _, oldRec := range oldRecs {
		newRec, err := updateRecordFromData(oldRec.Record, op.Data)
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
		err := no.ApplyNested(tx, root, root, oldRecs, clauses)
		if err != nil {
			return 0, err
		}
	}

	err := tx.Commit()
	if err != nil {
		return 0, err
	}
	return len(oldRecs), nil
}

func (op NestedUpdateManyOperation) ApplyNested(tx db.RWTx, root db.ModelRef, parent db.ModelRef, parents []*db.QueryResult, clauses []db.QueryClause) (err error) {
	cls, child := handleRelationshipWhere(tx, parent, op.Relationship, op.Where)
	clauses = append(clauses, cls...)
	q := tx.Query(root, clauses...)
	outs := q.All()

	for i, oldRec := range outs {
		newRec, err := updateRecordFromData(oldRec.Record, op.Data)
		if err != nil {
			return err
		}
		err = tx.Update(oldRec.Record, newRec)
		outs[i].Record = newRec
		if err != nil {
			return err
		}
	}
	for _, no := range op.Nested {
		err := no.ApplyNested(tx, root, child, outs, clauses)
		if err != nil {
			return err
		}
	}
	return
}
