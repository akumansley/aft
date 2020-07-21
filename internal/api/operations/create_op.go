package operations

import (
	"awans.org/aft/internal/db"
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

func (op CreateOperation) Apply(tx db.RWTx) (db.Record, error) {
	tx.Insert(op.Record)
	for _, no := range op.Nested {
		err := no.ApplyNested(tx, op.Record)
		if err != nil {
			return nil, err
		}
	}
	tx.Commit()
	return op.Record, nil
}

func (op NestedCreateOperation) ApplyNested(tx db.RWTx, parent db.Record) (err error) {
	tx.Insert(op.Record)
	tx.Connect(parent.ID(), op.Record.ID(), op.Relationship.ID())

	for _, no := range op.Nested {
		err = no.ApplyNested(tx, op.Record)
		if err != nil {
			return
		}
	}
	return nil
}

func findOneByID(tx db.Tx, modelID db.ID, id db.ID) (db.Record, error) {
	return tx.FindOne(modelID, db.EqID(id))
}

func (op NestedConnectOperation) ApplyNested(tx db.RWTx, parent db.Record) (err error) {
	rec, err := tx.FindOne(op.Relationship.Target().ID(), db.Eq(op.UniqueQuery.Key, op.UniqueQuery.Val))
	if err != nil {
		return
	}

	tx.Connect(parent.ID(), rec.ID(), op.Relationship.ID())
	return
}
