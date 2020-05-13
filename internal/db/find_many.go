package db

import (
	"awans.org/aft/internal/hold"
	"awans.org/aft/internal/model"
)

func (fc FieldCriterion) Matcher() hold.Matcher {
	return hold.Eq(fc.Key, fc.Val)
}

func (op FindManyOperation) Apply(tx Tx) []model.Record {
	return tx.FindMany(op.ModelName, op.Query)
}

func (tx holdTx) FindMany(modelName string, query Query) []model.Record {
	var matchers []hold.Matcher
	for _, fc := range query.FieldCriteria {
		matchers = append(matchers, fc.Matcher())
	}
	mi := tx.h.IterMatches(modelName, hold.And(matchers...))
	var hits []model.Record
	for val, ok := mi.Next(); ok; val, ok = mi.Next() {
		hits = append(hits, val)
	}
	return hits
}
