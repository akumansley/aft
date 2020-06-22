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
	Relationship db.Relationship
	Record       db.Record
	Nested       []NestedOperation
}

type NestedConnectOperation struct {
	Relationship db.Relationship
	UniqueQuery  UniqueQuery
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
	tx.Connect(parent, op.Record, op.Relationship)

	for _, no := range op.Nested {
		err = no.ApplyNested(tx, op.Record)
		if err != nil {
			return
		}
	}
	return nil
}

func findOneByID(tx db.Tx, modelID db.ModelID, id db.ID) (db.Record, error) {
	return tx.FindOne(modelID, db.EqID(id))
}

func (op NestedConnectOperation) ApplyNested(tx db.RWTx, parent db.Record) (err error) {
	rec, err := tx.FindOne(op.Relationship.Target.ID, db.Eq(op.UniqueQuery.Key, op.UniqueQuery.Val))
	if err != nil {
		return
	}

	tx.Connect(parent, rec, op.Relationship)
	return
}
