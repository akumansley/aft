package db

import (
	"awans.org/aft/internal/hold"
	"awans.org/aft/internal/model"
)

func (fc FieldCriterion) Matcher() hold.Matcher {
	return hold.Eq(fc.Key, fc.Val)
}

func (op FindManyOperation) Apply(tx Tx) []model.Record {
	var matchers []hold.Matcher
	for _, fc := range op.Query.FieldCriteria {
		matchers = append(matchers, fc.Matcher())
	}
	return tx.FindMany(op.ModelName, hold.And(matchers...))
}
