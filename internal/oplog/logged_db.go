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

func LoggedDB(l OpLog, d db.DB) db.DB {
	return loggedDB{inner: d, l: l}
}

func (l loggedDB) GetModel(modelName string) (model.Model, error) {
	return l.inner.GetModel(modelName)
}

func (l loggedDB) SaveModel(m model.Model) {
	l.inner.SaveModel(m)
}

func (l loggedDB) MakeStruct(s string) interface{} {
	return l.inner.MakeStruct(s)
}

func (l loggedDB) Insert(st interface{}) {
	// TODO log
	l.inner.Insert(st)
}

func (l loggedDB) Connect(from, to interface{}, fromRel model.Relationship) {
	// TODO log
	l.inner.Connect(from, to, fromRel)
}

func (l loggedDB) Resolve(st interface{}, inc db.Inclusion) {
	l.inner.Resolve(st, inc)
}

func (l loggedDB) FindOne(modelName string, uq db.UniqueQuery) (interface{}, error) {
	return l.inner.FindOne(modelName, uq)
}

func (l loggedDB) FindMany(modelName string, q db.Query) []interface{} {
	return l.inner.FindMany(modelName, q)
}
