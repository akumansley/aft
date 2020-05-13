package db

import (
	"awans.org/aft/er/q"
	"awans.org/aft/internal/model"
)

func (fc FieldCriterion) Matcher() q.Matcher {
	return q.Eq(fc.Key, fc.Val)
}

func (op FindManyOperation) Apply(tx Tx) []model.Record {
	return tx.FindMany(op.ModelName, op.Query)
}

func (tx holdTx) FindMany(modelName string, query Query) []model.Record {
	var matchers []q.Matcher
	for _, fc := range query.FieldCriteria {
		matchers = append(matchers, fc.Matcher())
	}
	mi := tx.h.IterMatches(modelName, q.And(matchers...))
	var hits []model.Record
	for val, ok := mi.Next(); ok; val, ok = mi.Next() {
		hits = append(hits, val)
	}
	return hits
}
