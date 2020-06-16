package api

import (
	"awans.org/aft/internal/db"
	"github.com/google/uuid"
)

type FieldCriterion struct {
	Key string
	Val interface{}
}

type AggregateRelationshipCriterion struct {
	RelationshipCriterion RelationshipCriterion
	Aggregation           db.Aggregation
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
	q := buildQuery(tx, op)
	qrs := q.All()
	results := []db.Record{}
	for _, qr := range qrs {
		results = append(results, qr.Record)
	}
	return results
}

func buildQuery(tx db.Tx, op FindManyOperation) db.Q {
	root := tx.Ref(op.ModelID)
	q := tx.Query(root)
	q = handleWhere(tx, q, root, op.Where)
	return q
}

func handleWhere(tx db.Tx, q db.Q, parent db.ModelRef, w Where) db.Q {
	for _, fc := range w.FieldCriteria {
		q = q.Filter(parent, fc.Matcher())
	}
	for _, rc := range w.RelationshipCriteria {
		q = handleRC(tx, q, parent, rc)
	}
	for _, arc := range w.AggregateRelationshipCriteria {
		q = handleARC(tx, q, parent, arc)
	}
	return q
}

func handleRC(tx db.Tx, q db.Q, parent db.ModelRef, rc RelationshipCriterion) db.Q {
	child := tx.Ref(rc.Binding.Dual().ModelID())
	on := parent.Rel(rc.Binding.Name())
	q = q.Join(child, on)
	q = handleWhere(tx, q, child, rc.Where)

	return q
}

func handleARC(tx db.Tx, q db.Q, parent db.ModelRef, arc AggregateRelationshipCriterion) db.Q {
	child := tx.Ref(arc.RelationshipCriterion.Binding.Dual().ModelID())
	on := parent.Rel(arc.RelationshipCriterion.Binding.Name())
	q = q.Join(child, on)
	q = q.Aggregate(child, arc.Aggregation)
	q = handleWhere(tx, q, child, arc.RelationshipCriterion.Where)

	return q
}
