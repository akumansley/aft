package operations

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

func newId(st db.Record) error {
	u := uuid.New()
	err := st.Set("id", u)
	return err
}

func (op CreateOperation) Apply(tx db.RWTx) (db.Record, error) {
	err := newId(op.Record)
	if err != nil {
		return nil, err
	}
	tx.Insert(op.Record)
	for _, no := range op.Nested {
		no.ApplyNested(tx, op.Record)
	}
	return op.Record, nil
}

func (op NestedCreateOperation) ApplyNested(tx db.RWTx, parent db.Record) (err error) {
	err = newId(op.Record)
	if err != nil {
		return err
	}
	tx.Insert(op.Record)
	tx.Connect(parent, op.Record, op.Relationship)
	return nil
}

func findOneById(tx db.Tx, modelName string, id uuid.UUID) (db.Record, error) {
	return tx.FindOne(modelName, "id", id)
}

func (op NestedConnectOperation) ApplyNested(tx db.RWTx, parent db.Record) (err error) {
	modelName := op.Relationship.TargetModel
	st, err := tx.FindOne(modelName, op.UniqueQuery.Key, op.UniqueQuery.Val)
	if err != nil {
		return
	}
	tx.Connect(parent, st, op.Relationship)
	return
}
