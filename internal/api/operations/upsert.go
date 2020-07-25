package operations

import (
	"awans.org/aft/internal/db"
	"fmt"
)

func (op UpsertOperation) Apply(tx db.RWTx) (*db.QueryResult, error) {
	fo := FindOneOperation{ModelID: op.ModelID, Where: op.Where}
	rec, err := fo.Apply(tx)
	if err != nil {
		return nil, err
	}
	//record not found, so create it
	if rec == nil {
		co := CreateOperation{
			Record:  op.Create,
			Include: op.Include,
			Nested:  op.NestedCreate,
		}
		return co.Apply(tx)
		//record was found, so update it
	} else {
		uo := UpdateOperation{
			ModelID: op.ModelID,
			Where:   op.Where,
			Data:    op.Update,
			Nested:  op.NestedUpdate,
			Include: op.Include,
		}
		return uo.Apply(tx)
	}
}

func (op NestedUpsertOperation) ApplyNested(tx db.RWTx) (err error) {
	parent := tx.Ref(op.Relationship.Source().ID())
	child := tx.Ref(op.Relationship.Target().ID())

	root := tx.Ref(op.Relationship.Target().ID())
	clauses := handleWhere(tx, root, op.Where)
	on := parent.Rel(op.Relationship)
	clauses = append(clauses, db.Join(child, on))
	q := tx.Query(root, clauses...)
	outs := q.All()

	if len(outs) == 1 {
		uo := NestedUpdateOperation{
			Relationship: op.Relationship,
			Data:         op.Update,
			Nested:       op.NestedUpdate,
		}
		return uo.ApplyNested(tx)
	} else if len(outs) == 0 {
		co := NestedCreateOperation{
			Relationship: op.Relationship,
			Data:         op.Create,
			Nested:       op.NestedCreate,
		}
		return co.ApplyNested(tx)
	}
	return fmt.Errorf("Found more than one record")
}
