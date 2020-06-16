package api

import (
	"awans.org/aft/internal/db"
	"github.com/google/uuid"
)

type NestedOperation interface {
	ApplyNested(db.RWTx, db.Record) error
}

type CreateOperation struct {
	Record db.Record
	Nested []NestedOperation
}

type NestedCreateOperation struct {
	Binding db.Binding
	Record  db.Record
	Nested  []NestedOperation
}

type NestedConnectOperation struct {
	Binding     db.Binding
	UniqueQuery UniqueQuery
}

type UniqueQuery struct {
	Key string
	Val interface{}
}

func newID(st db.Record) error {
	u := uuid.New()
	err := st.Set("id", u)
	return err
}

func (op CreateOperation) Apply(tx db.RWTx) (db.Record, error) {
	err := newID(op.Record)
	if err != nil {
		return nil, err
	}
	tx.Insert(op.Record)
	for _, no := range op.Nested {
		err = no.ApplyNested(tx, op.Record)
		if err != nil {
			return nil, err
		}
	}
	tx.Commit()
	return op.Record, nil
}

func (op NestedCreateOperation) ApplyNested(tx db.RWTx, parent db.Record) (err error) {
	err = newID(op.Record)
	if err != nil {
		return err
	}
	tx.Insert(op.Record)
	// the binding is the parent's side of the relationship
	if op.Binding.Left {
		tx.Connect(parent, op.Record, op.Binding.Relationship)
	} else {
		tx.Connect(op.Record, parent, op.Binding.Relationship)
	}
	for _, no := range op.Nested {
		err = no.ApplyNested(tx, op.Record)
		if err != nil {
			return
		}
	}
	return nil
}

func findOneByID(tx db.Tx, modelID uuid.UUID, id uuid.UUID) (db.Record, error) {
	return tx.FindOne(modelID, db.Eq("id", id))
}

func (op NestedConnectOperation) ApplyNested(tx db.RWTx, parent db.Record) (err error) {
	rec, err := tx.FindOne(op.Binding.Dual().ModelID(), db.Eq(op.UniqueQuery.Key, op.UniqueQuery.Val))
	if err != nil {
		return
	}

	if op.Binding.Left {
		tx.Connect(parent, rec, op.Binding.Relationship)
	} else {
		tx.Connect(rec, parent, op.Binding.Relationship)
	}
	return
}
