package operations

import (
	"awans.org/aft/internal/db"
)

func (op DeleteManyOperation) Apply(tx db.RWTx) (int, error) {
	root := tx.Ref(op.ModelID)
	clauses := HandleWhere(tx, root, op.Where)
	q := tx.Query(root, clauses...)
	outs := q.All()

	for _, no := range op.Nested {
		err := no.ApplyNested(tx, root, outs)
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

func (op NestedDeleteManyOperation) ApplyNested(tx db.RWTx, parent db.ModelRef, parents []*db.QueryResult) (err error) {
	outs, child := handleRelationshipWhere(tx, parent, parents, op.Relationship, op.Where)

	for _, no := range op.Nested {
		err := no.ApplyNested(tx, child, outs)
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
