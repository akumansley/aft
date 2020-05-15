package operations

import (
	"awans.org/aft/internal/db"
	"awans.org/aft/internal/hold"
	"awans.org/aft/internal/model"
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
	Relationship                         model.Relationship
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

func (fc FieldCriterion) Matcher() hold.Matcher {
	return hold.Eq(fc.Key, fc.Val)
}

func (op FindManyOperation) Apply(tx db.Tx) []model.Record {
	var matchers []hold.Matcher
	for _, fc := range op.Query.FieldCriteria {
		matchers = append(matchers, fc.Matcher())
	}
	return tx.FindMany(op.ModelName, hold.And(matchers...))
}
