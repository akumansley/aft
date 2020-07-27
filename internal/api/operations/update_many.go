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
	for _, no := range op.Nested {
		err := no.ApplyNested(tx)
		if err != nil {
			return 0, err
		}
	}
	err = tx.Commit()
	if err != nil {
		return 0, err
	}
	return len(oldRecs), nil
}

func (op NestedUpdateManyOperation) ApplyNested(tx db.RWTx) (err error) {
	parent := tx.Ref(op.Relationship.Source().ID())
	child := tx.Ref(op.Relationship.Target().ID())

	root := tx.Ref(op.Relationship.Target().ID())
	clauses := handleWhere(tx, root, op.Where)
	on := parent.Rel(op.Relationship)
	clauses = append(clauses, db.Join(child, on))
	q := tx.Query(root, clauses...)
	oldRecs := q.All()

	for _, oldRec := range oldRecs {
		newRec, err := updateRecordFromData(oldRec.Record, op.Data)
		if err != nil {
			return err
		}
		err = tx.Update(oldRec.Record, newRec)
		if err != nil {
			return err
		}
	}
	for _, no := range op.Nested {
		err := no.ApplyNested(tx)
		if err != nil {
			return err
		}
	}
	return
}
