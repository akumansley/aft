package oplog

import (
	"awans.org/aft/internal/db"
	"github.com/google/uuid"
)

type DBOp int

const (
	Connect DBOp = iota
	Create
	Update
	Delete
)

type TxEntry struct {
	Ops []DBOpEntry
}

func (txe TxEntry) Replay(rwtx db.RWTx) {
	for _, op := range txe.Ops {
		op.Replay(rwtx)
	}
}

type DBOpEntry struct {
	OpType DBOp
	Op     interface{}
}

func (oe DBOpEntry) Replay(rwtx db.RWTx) {
	switch oe.OpType {
	case Create:
		cro := oe.Op.(CreateOp)
		cro.Replay(rwtx)
	case Connect:
		cno := oe.Op.(ConnectOp)
		cno.Replay(rwtx)
	case Update:
		panic("Not implemented")
	case Delete:
		panic("Not implemented")
	}
}

type CreateOp struct {
	Record db.Record
}

func (cro CreateOp) Replay(rwtx db.RWTx) {
	rwtx.Insert(cro.Record)
}

type ConnectOp struct {
	Left  uuid.UUID
	Right uuid.UUID
	RelId uuid.UUID
}

func (cno ConnectOp) Replay(rwtx db.RWTx) {
	relRec, err := rwtx.FindOne("relationship", "id", cno.RelId)
	if err != nil {
		panic("couldn't find one on replay")
	}
	rel := db.LoadRel(relRec)
	leftModel, err := rwtx.GetModelById(rel.LeftModelId)
	if err != nil {
		panic("couldn't find one on replay")
	}
	left, err := rwtx.FindOne(leftModel.Name, "id", cno.Left)
	if err != nil {
		panic("couldn't find one on replay")
	}
	rightModel, err := rwtx.GetModelById(rel.RightModelId)
	if err != nil {
		panic("couldn't find one on replay")
	}
	right, err := rwtx.FindOne(rightModel.Name, "id", cno.Right)
	if err != nil {
		panic("couldn't find one on replay")
	}
	rwtx.Connect(left, right, rel)
}

type UpdateOp struct {
	Id     uuid.UUID
	Record db.Record
}

type DeleteOp struct {
	Id uuid.UUID
}

type loggedDB struct {
	inner db.DB
	l     OpLog
}

type loggedTx struct {
	inner db.RWTx
	txe   TxEntry
	l     OpLog
}

func DBFromLog(db db.DB, l OpLog) error {
	iter := l.Iterator()
	rwtx := db.NewRWTx()
	for iter.Next() {
		val := iter.Value()
		txe := val.(TxEntry)
		txe.Replay(rwtx)
	}
	if iter.Err() != nil {
		return iter.Err()
	}
	err := rwtx.Commit()
	return err
}

func LoggedDB(l OpLog, d db.DB) db.DB {
	return &loggedDB{inner: d, l: l}
}

func (l *loggedDB) NewTx() db.Tx {
	return l.inner.NewTx()
}

func (l *loggedDB) NewRWTx() db.RWTx {
	return &loggedTx{inner: l.inner.NewRWTx(), l: l.l}
}

func (l *loggedDB) DeepEquals(o db.DB) bool {
	return l.inner.DeepEquals(o)
}

func (l *loggedDB) Iterator() db.Iterator {
	return l.inner.Iterator()
}

func (tx *loggedTx) GetModelById(id uuid.UUID) (db.Model, error) {
	return tx.inner.GetModelById(id)
}

func (tx *loggedTx) GetModel(modelName string) (db.Model, error) {
	return tx.inner.GetModel(modelName)
}

func (tx *loggedTx) SaveModel(m db.Model) {
	tx.inner.SaveModel(m)
}

func (tx *loggedTx) MakeRecord(s string) db.Record {
	return tx.inner.MakeRecord(s)
}

func (tx *loggedTx) Insert(rec db.Record) {
	co := CreateOp{Record: rec}
	dboe := DBOpEntry{Create, co}
	tx.txe.Ops = append(tx.txe.Ops, dboe)
	tx.inner.Insert(rec)
}

func (tx *loggedTx) Connect(left, right db.Record, rel db.Relationship) {
	co := ConnectOp{Left: left.Id(), Right: right.Id(), RelId: rel.Id}
	dboe := DBOpEntry{Connect, co}
	tx.txe.Ops = append(tx.txe.Ops, dboe)
	tx.inner.Connect(left, right, rel)
}

func (tx *loggedTx) FindOne(modelName string, key string, val interface{}) (db.Record, error) {
	return tx.inner.FindOne(modelName, key, val)
}

func (tx *loggedTx) FindMany(modelName string, matcher db.Matcher) []db.Record {
	return tx.inner.FindMany(modelName, matcher)
}

func (tx *loggedTx) Commit() (err error) {
	err = tx.l.Log(tx.txe)
	if err != nil {
		return
	}
	err = tx.inner.Commit()
	return
}
