package operations

import (
	"awans.org/aft/internal/db"
)

type UniqueQuery struct {
	Key string
	Val interface{}
}

func (op CreateOperation) Apply(tx db.RWTx) (*db.QueryResult, error) {
	rec, err := buildRecordFromData(tx, op.ModelID, op.Data)
	if err != nil {
		return nil, err
	}
	tx.Insert(rec)

	root := tx.Ref(rec.Interface().ID())
	parents := []*db.QueryResult{&db.QueryResult{Record: rec}}
	for _, no := range op.Nested {
		err := no.ApplyNested(tx, op.Record)
		if err != nil {
			return nil, err
		}
	}
	out, err := op.FindManyArgs.Include.One(tx, op.Record.Interface().ID(), op.Record)
	if err != nil {
		return nil, err
	}
	tx.Commit()
	return qrs[0], nil
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

func (op NestedConnectOperation) ApplyNested(tx db.RWTx, parent db.Record) (err error) {
	t := tx.Ref(op.Relationship.Target().ID())
	res, err := tx.Query(t).Filter(t, db.Eq(op.UniqueQuery.Key, op.UniqueQuery.Val)).One()
	rec := res.Record
	if err != nil {
		return
	}

	tx.Connect(parent.ID(), rec.ID(), op.Relationship.ID())
	return
}

func buildRecordFromData(tx db.RWTx, modelID db.ID, data map[string]interface{}) (db.Record, error) {
	m, err := tx.Schema().GetInterfaceByID(modelID)
	if err != nil {
		return nil, err
	}
	rec := db.NewRecord(m)
	for k, v := range data {
		rec.Set(k, v)
	}
	return rec, nil
}
