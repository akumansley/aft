package operations

import (
	"awans.org/aft/internal/db"
)

func (op DeleteManyOperation) Apply(tx db.RWTx) (int, error) {
	root := tx.Ref(op.ModelID)
	clauses := handleWhere(tx, root, op.Where)
	q := tx.Query(root, clauses...)
	outs := q.All()

	for _, no := range op.Nested {
		err := no.ApplyNested(tx, root, root, outs, clauses)
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

func (op NestedDeleteManyOperation) ApplyNested(tx db.RWTx, root db.ModelRef, parent db.ModelRef, parents []*db.QueryResult, clauses []db.QueryClause) (err error) {
	cls, child := handleRelationshipWhere(tx, parent, op.Relationship, op.Where)
	clauses = append(clauses, cls...)
	q := tx.Query(root, clauses...)
	outs := getEdgeResults(parents, q.All())

	for _, no := range op.Nested {
		err := no.ApplyNested(tx, root, child, outs, clauses)
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
