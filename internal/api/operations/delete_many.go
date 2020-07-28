package operations

import (
	"awans.org/aft/internal/db"
)

func (op DeleteManyOperation) Apply(tx db.RWTx) (int, error) {
	fm := FindManyOperation{ModelID: op.ModelID, FindManyArgs: FindManyArgs{Where: op.Where}}
	outs, err := fm.Apply(tx)
	if err != nil {
		return 0, err
	}
	for _, no := range op.Nested {
		err := no.ApplyNested(tx)
		if err != nil {
			return 0, err
		}
	}
	for _, out := range outs {
		err := tx.Delete(out.Record)
		if err != nil {
			return 0, err
		}
	}
	tx.Commit()
	return len(outs), nil
}

func (op NestedDeleteManyOperation) ApplyNested(tx db.RWTx) (err error) {
	parent := tx.Ref(op.Relationship.Source().ID())
	child := tx.Ref(op.Relationship.Target().ID())

	root := tx.Ref(op.Relationship.Target().ID())
	clauses := handleWhere(tx, root, op.Where)
	on := parent.Rel(op.Relationship)
	clauses = append(clauses, db.Join(child, on))
	q := tx.Query(root, clauses...)
	outs := q.All()

	for _, no := range op.Nested {
		err := no.ApplyNested(tx)
		if err != nil {
			return err
		}
	}
	for _, out := range outs {
		err = tx.Delete(out.Record)
		return err
	}
	return
}
