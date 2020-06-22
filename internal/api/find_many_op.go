package api

import (
	"awans.org/aft/internal/db"
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
	Relationship db.Relationship
	Where        Where
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
	ModelID db.ModelID
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
	qb := q.AsBlock()
	qb = handleWhere(tx, qb, root, op.Where)
	q.SetMainBlock(qb)
	return q
}

func handleWhere(tx db.Tx, q db.QBlock, parent db.ModelRef, w Where) db.QBlock {
	for _, fc := range w.FieldCriteria {
		q = q.Filter(parent, fc.Matcher())
	}
	for _, rc := range w.RelationshipCriteria {
		q = handleRC(tx, q, parent, rc)
	}
	for _, arc := range w.AggregateRelationshipCriteria {
		q = handleARC(tx, q, parent, arc)
	}

	var orBlocks []db.QBlock
	for _, or := range w.Or {
		orBlock := handleSetOpBranch(tx, parent, or)
		orBlocks = append(orBlocks, orBlock)
	}
	if len(orBlocks) > 0 {
		q = q.Or(parent, orBlocks...)
	}

	var andBlocks []db.QBlock
	for _, and := range w.And {
		andBlock := handleSetOpBranch(tx, parent, and)
		andBlocks = append(andBlocks, andBlock)
	}
	if len(andBlocks) > 0 {
		q = q.And(parent, andBlocks...)
	}

	var notBlocks []db.QBlock
	for _, not := range w.Not {
		notBlock := handleSetOpBranch(tx, parent, not)
		notBlocks = append(notBlocks, notBlock)
	}
	if len(notBlocks) > 0 {
		q = q.Not(parent, notBlocks...)
	}
	return q
}

func handleSetOpBranch(tx db.Tx, parent db.ModelRef, w Where) db.QBlock {
	qb := db.NewBlock()
	return handleWhere(tx, qb, parent, w)
}

func handleRC(tx db.Tx, q db.QBlock, parent db.ModelRef, rc RelationshipCriterion) db.QBlock {
	child := tx.Ref(rc.Relationship.Target.ID)
	on := parent.Rel(rc.Relationship)
	q = q.Join(child, on)
	q = handleWhere(tx, q, child, rc.Where)

	return q
}

func handleARC(tx db.Tx, q db.QBlock, parent db.ModelRef, arc AggregateRelationshipCriterion) db.QBlock {
	child := tx.Ref(arc.RelationshipCriterion.Relationship.Target.ID)
	on := parent.Rel(arc.RelationshipCriterion.Relationship)

	q = q.Join(child, on)
	q = q.Aggregate(child, arc.Aggregation)
	q = handleWhere(tx, q, child, arc.RelationshipCriterion.Where)

	return q
}
