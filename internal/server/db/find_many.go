package db

import (
	"awans.org/aft/er/q"
	"fmt"
)

// I guess this is a needless copy
func (fc FieldCriterion) Matcher() q.Matcher {
	return q.Eq(fc.Key, fc.Val)
}

func (op FindManyOperation) Apply(db DB) interface{} {
	var matchers []q.Matcher
	for _, fc := range op.Query.FieldCriteria {
		matchers = append(matchers, fc.Matcher())
	}
	mi := db.h.IterMatches(op.ModelName, q.And(matchers...))
	var hits []interface{}
	for val, ok := mi.Next(); ok; val, ok = mi.Next() {
		hits = append(hits, val)
	}
	return hits
}
