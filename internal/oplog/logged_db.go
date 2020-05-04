package oplog

import (
	"awans.org/aft/internal/model"
	"awans.org/aft/internal/server/db"
)

type DBOp struct {
	st interface{}
	Op int
}

type loggedDB struct {
	inner db.DB
	l     OpLog
}

type loggedTx struct {
	inner db.RWTx
}

func LoggedDB(l OpLog, d db.DB) db.DB {
	return &loggedDB{inner: d, l: l}
}

func (l *loggedDB) NewTx() db.Tx {
	return l.inner.NewTx()
}

func (l *loggedDB) NewRWTx() db.RWTx {
	return &loggedTx{inner: l.inner.NewRWTx()}
}

func (tx *loggedTx) GetModel(modelName string) (model.Model, error) {
	return tx.inner.GetModel(modelName)
}

func (tx *loggedTx) SaveModel(m model.Model) {
	tx.inner.SaveModel(m)
}

func (tx *loggedTx) MakeStruct(s string) interface{} {
	return tx.inner.MakeStruct(s)
}

func (tx *loggedTx) Insert(st interface{}) {
	// TODO log
	tx.inner.Insert(st)
}

func (tx *loggedTx) Connect(from, to interface{}, fromRel model.Relationship) {
	// TODO log
	tx.inner.Connect(from, to, fromRel)
}

func (tx *loggedTx) Resolve(st interface{}, inc db.Inclusion) {
	tx.inner.Resolve(st, inc)
}

func (tx *loggedTx) FindOne(modelName string, uq db.UniqueQuery) (interface{}, error) {
	return tx.inner.FindOne(modelName, uq)
}

func (tx *loggedTx) FindMany(modelName string, q db.Query) []interface{} {
	return tx.inner.FindMany(modelName, q)
}

func (tx *loggedTx) Commit() {
	tx.inner.Commit()
}
