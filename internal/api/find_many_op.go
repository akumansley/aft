package api

import (
	"awans.org/aft/internal/db"
)

type FieldCriterion struct {
	Key string
	Val interface{}
}

type Aggregation int

const (
	Every Aggregation = iota
	Some
	None
	Single
)

type AggregateRelationshipCriterion struct {
	RelationshipCriterion RelationshipCriterion
	Aggregation           Aggregation
}

type RelationshipCriterion struct {
	Binding                              db.Binding
	RelatedFieldCriteria                 []FieldCriterion
	RelatedRelationshipCriteria          []RelationshipCriterion
	RelatedAggregateRelationshipCriteria []AggregateRelationshipCriterion
}

type Query struct {
	FieldCriteria                 []FieldCriterion
	RelationshipCriteria          []RelationshipCriterion
	AggregateRelationshipCriteria []AggregateRelationshipCriterion
	Or                            []Query
	And                           []Query
	Not                           []Query
}

type FindManyOperation struct {
	ModelName string
	Query     Query
}

func (fc FieldCriterion) Matcher() db.Matcher {
	return db.Eq(fc.Key, fc.Val)
}

func (op FindManyOperation) Apply(tx db.Tx) []db.Record {
	var matchers []db.Matcher
	for _, fc := range op.Query.FieldCriteria {
		matchers = append(matchers, fc.Matcher())
	}
	return tx.FindMany(op.ModelName, db.And(matchers...))
}
