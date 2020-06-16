package api

import (
	"awans.org/aft/internal/db"
)

// type sort struct {
// }
// type limit struct {
// 	r     relation
// 	limit int
// }

// type setop int

// const (
// 	union setop = iota
// 	subtraction
// 	intersection
// )

// type hashsetop struct {
// 	left, right relation
// 	lsid, rsid  int
// 	op          setop
// }

// func (s *hashsetop) String() string {
// 	return fmt.Sprintf("hashsetop{%v, %v, %v}", s.left, s.right, s.op)
// }

// func (s *hashsetop) iter(tx db.Tx) iterator {
// }

func Plan(fm FindManyOperation) *db.Query {
	m := db.Ref(fm.ModelID)
	q := db.Select(m)
	planWhere(m, q, fm.Where)
	return q
}

func planWhere(m db.ModelRef, q *db.Query, w Where) {

	// first we filter locally
	for _, fc := range w.FieldCriteria {
		q.Where(m.Match(fc.Matcher()))
	}

	// then we join to-one related fields
	for _, rc := range w.RelationshipCriteria {
		planRC(rc, m, q)
	}

	// then we join to-many related fields
	for _, arc := range w.AggregateRelationshipCriteria {
		planARC(arc, m, q)
	}

	// 	// then we union by or
	// 	for _, orWhere := range w.Or {
	// 		orSid, orR := planWhere(modelID, orWhere, jt)
	// 		r = &hashsetop{r, orR, sid, orSid, union}
	// 	}

	// 	// then we subtract not
	// 	// should this be conceptualized as an "antijoin" in the typical case
	// 	for _, notWhere := range w.Not {
	// 		notSid, notR := planWhere(modelID, notWhere, jt)
	// 		r = &hashsetop{r, notR, sid, notSid, subtraction}
	// 	}

	// 	// then we intersect and
	// 	for _, andWhere := range w.And {
	// 		andSid, andR := planWhere(modelID, andWhere, jt)
	// 		r = &hashsetop{r, andR, sid, andSid, intersection}
	// 	}
}

func planRC(rc RelationshipCriterion, parent db.ModelRef, q *db.Query) {
	d := rc.Binding.Dual()
	child := db.Ref(d.ModelID())
	q.JoinOne(db.Select(child), db.On(parent, child, rc.Binding))
	planWhere(child, q, rc.Where)
}

func planARC(arc AggregateRelationshipCriterion, parent db.ModelRef, q *db.Query) {
	b := arc.RelationshipCriterion.Binding
	d := b.Dual()

	child := db.Ref(d.ModelID())
	q.JoinMany(db.Select(child), db.On(parent, child, b), arc.Aggregation)
	planWhere(child, q, arc.RelationshipCriterion.Where)
}
