package operations

import (
	"awans.org/aft/internal/db"
	"fmt"
)

func (op DeleteOperation) Apply(tx db.RWTx) (*db.QueryResult, error) {
	root := tx.Ref(op.ModelID)
	clauses := handleFindMany(tx, root, op.FindArgs)
	q := tx.Query(root, clauses...)
	outs := q.All()

	if len(outs) > 1 {
		return nil, fmt.Errorf("Found more than one record")
	}

	for _, no := range op.Nested {
		err := no.ApplyNested(tx, root, root, outs, clauses)
		if err != nil {
			return nil, err
		}
	}

	if len(outs) == 0 {
		return nil, nil
	}
	out := outs[0]
	err := tx.Delete(out.Record)
	if err != nil {
		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}
	return out, err
}

func (op NestedDeleteOperation) ApplyNested(tx db.RWTx, root db.ModelRef, parent db.ModelRef, parents []*db.QueryResult, clauses []db.QueryClause) (err error) {
	cls, child := handleRelationshipWhere(tx, parent, op.Relationship, op.Where)
	clauses = append(clauses, cls...)
	q := tx.Query(root, clauses...)
	outs := getEdgeResults(parents, q.All())
	if len(outs) > 1 {
		return fmt.Errorf("Found more than one record")
	} else if len(outs) == 1 {
		for _, no := range op.Nested {
			err := no.ApplyNested(tx, root, child, outs, clauses)
			if err != nil {
				return err
			}
		}
		tx.Delete(outs[0].Record)
	}
	return
}
