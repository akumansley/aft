package operations

import (
	"awans.org/aft/internal/db"
	"fmt"
)

func (op UpsertOperation) Apply(tx db.RWTx) (*db.QueryResult, error) {
	root := tx.Ref(op.ModelID)
	clauses := handleFindMany(tx, root, op.FindArgs)
	q := tx.Query(root, clauses...)
	outs := q.All()

	if len(outs) > 1 {
		return nil, fmt.Errorf("Found more than one record")
	}
	if len(outs) == 0 {
		co := CreateOperation{
			Record:   op.Create,
			FindArgs: op.FindArgs,
			Nested:   op.NestedCreate,
		}
		return co.Apply(tx)
	} else {
		uo := UpdateOperation{
			ModelID:  op.ModelID,
			FindArgs: op.FindArgs,
			Data:     op.Update,
			Nested:   op.NestedUpdate,
		}
		return uo.Apply(tx)
	}
}

func (op NestedUpsertOperation) ApplyNested(tx db.RWTx, root db.ModelRef, parent db.ModelRef, parents []*db.QueryResult, clauses []db.QueryClause) (err error) {
	cls, child := handleRelationshipWhere(tx, parent, op.Relationship, op.Where)
	clauses = append(clauses, cls...)
	q := tx.Query(root, clauses...)
	outs := getEdgeResults(parents, q.All())

	if len(outs) == 1 {
		uo := NestedUpdateOperation{
			Relationship: op.Relationship,
			Data:         op.Update,
			Nested:       op.NestedUpdate,
		}
		return uo.ApplyNested(tx, root, child, parents, clauses)
	} else if len(outs) == 0 {
		co := NestedCreateOperation{
			Relationship: op.Relationship,
			Record:       op.Create,
			Nested:       op.NestedCreate,
		}
		return co.ApplyNested(tx, root, child, parents, clauses)
	}
	return fmt.Errorf("Found more than one record")
}
