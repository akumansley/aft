package operations

import (
	"awans.org/aft/internal/db"
	"awans.org/aft/internal/model"
	"github.com/google/uuid"
)

type NestedOperation interface {
	ApplyNested(db.RWTx, model.Record) error
}

type CreateOperation struct {
	Record model.Record
	Nested []NestedOperation
}

type NestedCreateOperation struct {
	Relationship model.Relationship
	Record       model.Record
	Nested       []NestedOperation
}

type NestedConnectOperation struct {
	Relationship model.Relationship
	UniqueQuery  UniqueQuery
}

type UniqueQuery struct {
	Key string
	Val interface{}
}

func newId(st model.Record) {
	u := uuid.New()
	st.Set("id", u)
}

func (op CreateOperation) Apply(tx db.RWTx) (model.Record, error) {
	newId(op.Record)
	tx.Insert(op.Record)
	for _, no := range op.Nested {
		no.ApplyNested(tx, op.Record)
	}
	return op.Record, nil
}

func (op NestedCreateOperation) ApplyNested(tx db.RWTx, parent model.Record) (err error) {
	newId(op.Record)
	tx.Insert(op.Record)
	tx.Connect(parent, op.Record, op.Relationship)
	return nil
}

func findOneById(tx db.Tx, modelName string, id uuid.UUID) (model.Record, error) {
	return tx.FindOne(modelName, "id", id)
}

func (op NestedConnectOperation) ApplyNested(tx db.RWTx, parent model.Record) (err error) {
	modelName := op.Relationship.TargetModel
	st, err := tx.FindOne(modelName, op.UniqueQuery.Key, op.UniqueQuery.Val)
	if err != nil {
		return
	}
	tx.Connect(parent, st, op.Relationship)
	return
}
