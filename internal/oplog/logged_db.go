package oplog

import (
	"errors"

	"awans.org/aft/internal/db"
)

func DBFromLog(appDB db.DB, l OpLog) error {
	iter := l.Iterator()
	for iter.Next() {
		rwtx := appDB.NewRWTx()
		val := iter.Value()
		logOps := val.([]LogOp)
		var ops []db.Operation
		for _, logOp := range logOps {
			op, err := logOp.Decode(rwtx)
			if err != nil {
				return err
			}
			ops = append(ops, op)
		}
		for _, op := range ops {
			op.Replay(rwtx)
		}
		err := rwtx.Commit()
		if err != nil {
			return err
		}
	}
	if iter.Err() != nil {
		return iter.Err()
	}
	return nil
}

func modelForAttr(tx db.Tx, attrID db.ID) db.Model {
	models := tx.Ref(db.ModelModel.ID())
	attrs := tx.Ref(db.ConcreteAttributeModel.ID())
	q := tx.Query(models,
		db.Join(attrs, models.Rel(db.ModelAttributes)),
		db.Aggregate(attrs, db.Some),
		db.Filter(attrs, db.EqID(attrID)))
	rec, err := q.OneRecord()
	if err != nil {
		panic(err)
	}
	model := tx.Schema().LoadModel(rec)
	return model
}

func MakeTransactionLogger(l OpLog) func(db.BeforeCommit) {
	logger := func(event db.BeforeCommit) {
		ops := event.Tx.Operations()

		if len(ops) == 0 {
			return
		}
		encoded := encode(ops)
		err := l.Log(encoded)
		if err != nil {
			panic(err)
		}
	}
	return logger
}

func encode(ops []db.Operation) []LogOp {
	var encoded []LogOp

	for _, op := range ops {
		switch op.(type) {
		case db.CreateOp:
			c := op.(db.CreateOp)
			enc := create{c.Record.RawData(), c.ModelID}
			encoded = append(encoded, enc)
		case db.DeleteOp:
			d := op.(db.DeleteOp)
			enc := delete{d.Record.RawData(), d.ModelID}
			encoded = append(encoded, enc)
		case db.UpdateOp:
			u := op.(db.UpdateOp)
			enc := update{u.OldRecord.RawData(), u.NewRecord.RawData(), u.ModelID}
			encoded = append(encoded, enc)
		case db.ConnectOp:
			c := op.(db.ConnectOp)
			enc := connect{c.Source, c.Target, c.RelID}
			encoded = append(encoded, enc)
		case db.DisconnectOp:
			d := op.(db.DisconnectOp)
			enc := disconnect{d.Source, d.Target, d.RelID}
			encoded = append(encoded, enc)
		default:
			panic(errors.New("invalid op"))
		}
	}
	return encoded
}

type LogOp interface {
	Decode(db.Tx) (db.Operation, error)
}

type create struct {
	RecData interface{}
	ModelID db.ID
}

func (c create) Decode(tx db.Tx) (db.Operation, error) {
	model, err := tx.Schema().GetModelByID(c.ModelID)
	if err != nil {
		return nil, err
	}
	return db.CreateOp{
		Record:  db.RecordFromParts(c.RecData, model),
		ModelID: c.ModelID,
	}, nil
}

type update struct {
	OldRecData interface{}
	NewRecData interface{}
	ModelID    db.ID
}

func (u update) Decode(tx db.Tx) (db.Operation, error) {
	model, err := tx.Schema().GetModelByID(u.ModelID)
	if err != nil {
		return nil, err
	}
	return db.UpdateOp{
		OldRecord: db.RecordFromParts(u.OldRecData, model),
		NewRecord: db.RecordFromParts(u.NewRecData, model),
		ModelID:   u.ModelID,
	}, nil
}

type delete struct {
	RecData interface{}
	ModelID db.ID
}

func (d delete) Decode(tx db.Tx) (db.Operation, error) {
	model, err := tx.Schema().GetModelByID(d.ModelID)
	if err != nil {
		return nil, err
	}
	return db.DeleteOp{
		Record:  db.RecordFromParts(d.RecData, model),
		ModelID: d.ModelID,
	}, nil
}

type connect struct {
	Source db.ID
	Target db.ID
	RelID  db.ID
}

func (c connect) Decode(tx db.Tx) (db.Operation, error) {
	return db.ConnectOp{
		Source: c.Source, Target: c.Target, RelID: c.RelID,
	}, nil
}

type disconnect struct {
	Source db.ID
	Target db.ID
	RelID  db.ID
}

func (d disconnect) Decode(tx db.Tx) (db.Operation, error) {
	return db.DisconnectOp{
		Source: d.Source, Target: d.Target, RelID: d.RelID,
	}, nil
}
