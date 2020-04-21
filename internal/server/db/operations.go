package db

import (
	"awans.org/aft/internal/model"
)

type Operation interface {
	Apply(DB) interface{}
}

type NestedOperation interface {
	ApplyNested(DB, interface{})
}

type CreateOperation struct {
	Struct interface{}
	Nested []NestedOperation
}

type NestedCreateOperation struct {
	Relationship model.Relationship
	Struct       interface{}
	Nested       []NestedOperation
}

type UniqueQuery struct {
	Key string
	Val interface{}
}

type NestedConnectOperation struct {
	Relationship model.Relationship
	UniqueQuery  UniqueQuery
}

type FindOneOperation struct {
	ModelName   string
	UniqueQuery UniqueQuery
}

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
