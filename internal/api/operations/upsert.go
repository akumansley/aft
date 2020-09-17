package operations

import (
	"fmt"

	"awans.org/aft/internal/db"
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
			ModelID:  op.ModelID,
			Data:     op.Create,
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

func (op NestedUpsertOperation) ApplyNested(tx db.RWTx, parent db.ModelRef, parents []*db.QueryResult) (err error) {
	outs, _ := handleRelationshipWhere(tx, parent, parents, op.Relationship, op.Where)

	if len(outs) == 1 {
		uo := NestedUpdateOperation{
			Relationship: op.Relationship,
			Data:         op.Update,
			Nested:       op.NestedUpdate,
			Where:        op.Where,
		}
		return uo.ApplyNested(tx, parent, parents)
	} else if len(outs) == 0 {
		co := NestedCreateOperation{
			Relationship: op.Relationship,
			Data:         op.Create,
			Nested:       op.NestedCreate,
		}
		return co.ApplyNested(tx, parent, parents)
	}
	return fmt.Errorf("Found more than one record")
}
