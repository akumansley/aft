package oplog

import (
	"awans.org/aft/internal/model"
	"awans.org/aft/internal/server/db"
	"github.com/google/uuid"
	"github.com/ompluscator/dynamic-struct"
)

func getId(st interface{}) uuid.UUID {
	reader := dynamicstruct.NewReader(st)
	id := reader.GetField("Id").Interface().(uuid.UUID)
	return id
}

type DBOp int

const (
	Create DBOp = iota
	Connect
	Update
	Delete
)

type TxEntry struct {
	ops []DBOpEntry
}

type DBOpEntry struct {
	OpType DBOp
	Op     interface{}
}

type CreateOp struct {
	Record model.Record
}

type ConnectOp struct {
	From  uuid.UUID
	To    uuid.UUID
	RelId uuid.UUID
}

type UpdateOp struct {
	Id     uuid.UUID
	Record model.Record
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

// func DBFromLog(db db.DB, l OpLog) db.DB {
// 	iter := l.Iterator()
// 	for op, ok := iter.Next(); ok; op, ok := iter.Next() {

// 	}
// }

func LoggedDB(l OpLog, d db.DB) db.DB {
	return &loggedDB{inner: d, l: l}
}

func (l *loggedDB) NewTx() db.Tx {
	return l.inner.NewTx()
}

func (l *loggedDB) NewRWTx() db.RWTx {
	return &loggedTx{inner: l.inner.NewRWTx(), l: l.l}
}

func (tx *loggedTx) GetModel(modelName string) (model.Model, error) {
	return tx.inner.GetModel(modelName)
}

func (tx *loggedTx) SaveModel(m model.Model) {
	tx.inner.SaveModel(m)
}

func (tx *loggedTx) MakeRecord(s string) model.Record {
	return tx.inner.MakeRecord(s)
}

func (tx *loggedTx) Insert(rec model.Record) {
	co := CreateOp{Record: rec}
	dboe := DBOpEntry{Create, co}
	tx.txe.ops = append(tx.txe.ops, dboe)
	tx.inner.Insert(rec)
}

func (tx *loggedTx) Connect(from, to model.Record, fromRel model.Relationship) {
	co := ConnectOp{From: from.Id(), To: to.Id()}
	dboe := DBOpEntry{Connect, co}
	tx.txe.ops = append(tx.txe.ops, dboe)
	tx.inner.Connect(from, to, fromRel)
}

func (tx *loggedTx) Resolve(ir *model.IncludeResult, inc db.Inclusion) {
	tx.inner.Resolve(ir, inc)
}

func (tx *loggedTx) FindOne(modelName string, uq db.UniqueQuery) (model.Record, error) {
	return tx.inner.FindOne(modelName, uq)
}

func (tx *loggedTx) FindMany(modelName string, q db.Query) []model.Record {
	return tx.inner.FindMany(modelName, q)
}

func (tx *loggedTx) Commit() {
	tx.l.Log(tx.txe)
	tx.inner.Commit()
}
