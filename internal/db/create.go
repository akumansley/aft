package db

import (
	"awans.org/aft/internal/model"
	"github.com/google/uuid"
)

func newId(st model.Record) {
	u := uuid.New()
	st.Set("id", u)
}

func (op CreateOperation) Apply(tx RWTx) (model.Record, error) {
	newId(op.Record)
	tx.Insert(op.Record)
	for _, no := range op.Nested {
		no.ApplyNested(tx, op.Record)
	}
	return op.Record, nil
}

func (op NestedCreateOperation) ApplyNested(tx RWTx, parent model.Record) (err error) {
	newId(op.Record)
	tx.Insert(op.Record)
	tx.Connect(parent, op.Record, op.Relationship)
	return nil
}

func findOneById(tx Tx, modelName string, id uuid.UUID) (model.Record, error) {
	return tx.FindOne(modelName, UniqueQuery{Key: "Id", Val: id})
}

func (op NestedConnectOperation) ApplyNested(tx RWTx, parent model.Record) (err error) {
	modelName := op.Relationship.TargetModel
	st, err := tx.FindOne(modelName, op.UniqueQuery)
	if err != nil {
		return
	}
	tx.Connect(parent, st, op.Relationship)
	return
}
