package db

import (
	"awans.org/aft/er/q"
)

// I guess this is a needless copy
func (fc FieldCriterion) Matcher() q.Matcher {
	return q.Eq(fc.Key, fc.Val)
}

func (op FindManyOperation) Apply(db DB) interface{} {
	return db.FindMany(op.ModelName, op.Query)
}

func (db holdDB) FindMany(modelName string, query Query) []interface{} {
	var matchers []q.Matcher
	for _, fc := range query.FieldCriteria {
		matchers = append(matchers, fc.Matcher())
	}
	mi := db.h.IterMatches(modelName, q.And(matchers...))
	var hits []interface{}
	for val, ok := mi.Next(); ok; val, ok = mi.Next() {
		hits = append(hits, val)
	}
	return hits
}
