package operations

import (
	"awans.org/aft/internal/db"
	"fmt"
)

func (op DeleteOperation) Apply(tx db.RWTx) (*db.QueryResult, error) {
	fo := FindOneOperation{ModelID: op.ModelID, FindManyArgs: op.FindManyArgs}
	out, err := fo.Apply(tx)
	if err != nil {
		return nil, err
	}
	if out == nil {
		return nil, fmt.Errorf("Didn't find record to delete")
	}
	inc, err := op.FindManyArgs.Include.One(tx, out.Record.Interface().ID(), out.Record)

	for _, no := range op.Nested {
		err := no.ApplyNested(tx)
		if err != nil {
			return nil, err
		}
	}
	err = tx.Delete(out.Record)
	if err != nil {
		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}
	return inc, err
}

func (op NestedDeleteOperation) ApplyNested(tx db.RWTx) (err error) {
	parent := tx.Ref(op.Relationship.Source().ID())
	child := tx.Ref(op.Relationship.Target().ID())

	root := tx.Ref(op.Relationship.Target().ID())
	clauses := handleWhere(tx, root, op.Where)
	on := parent.Rel(op.Relationship)
	clauses = append(clauses, db.Join(child, on))
	q := tx.Query(root, clauses...)
	out := q.All()

	if len(out) > 1 {
		return fmt.Errorf("Found more than one record")
	} else if len(out) == 1 {
		for _, no := range op.Nested {
			err := no.ApplyNested(tx)
			if err != nil {
				return err
			}
		}
		tx.Delete(out[0].Record)
	}
	return
}
