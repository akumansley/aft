package api

import (
	"awans.org/aft/internal/db"
	"github.com/google/uuid"
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
	Binding db.Binding
	Where   Where
}

type Where struct {
	FieldCriteria                 []FieldCriterion
	RelationshipCriteria          []RelationshipCriterion
	AggregateRelationshipCriteria []AggregateRelationshipCriterion
	Or                            []Where
	And                           []Where
	Not                           []Where
}

type FindManyOperation struct {
	ModelID uuid.UUID
	Where   Where
}

func (fc FieldCriterion) Matcher() db.Matcher {
	return db.Eq(fc.Key, fc.Val)
}

func (op FindManyOperation) Apply(tx db.Tx) []db.Record {
	var matchers []db.Matcher
	for _, fc := range op.Where.FieldCriteria {
		matchers = append(matchers, fc.Matcher())
	}
	return tx.FindMany(op.ModelID, db.And(matchers...))
}
