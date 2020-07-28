package operations

import (
	"awans.org/aft/internal/db"
)

func (op UpdateOperation) Apply(tx db.RWTx) (*db.QueryResult, error) {
	err := tx.Update(op.Old, op.New)
	if err != nil {
		return nil, err
	}
	for _, no := range op.Nested {
		err := no.ApplyNested(tx, op.Old)
		if err != nil {
			return nil, err
		}
	}
	err = tx.Commit()
	if err != nil {
		return nil, err
	}
	inc, err := op.Include.One(tx, op.Old.Interface().ID(), op.New)
	return inc, err
}

func (op NestedUpdateOperation) ApplyNested(tx db.RWTx, parent db.Record) (err error) {
	err = tx.Update(op.Old, op.New)
	if err != nil {
		return
	}
	for _, no := range op.Nested {
		err := no.ApplyNested(tx, op.Old)
		if err != nil {
			return err
		}
	}
	err = tx.Commit()
	if err != nil {
		return err
	}
	return err
}
