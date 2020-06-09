package api

import (
	"awans.org/aft/internal/db"
	"fmt"
	"github.com/google/uuid"
)

type relation interface {
	iter()
	String() string
}

type scan interface {
	iter()
	filter(db.Matcher)
	String() string
}

type seqscan struct {
	modelID  uuid.UUID
	matchers []db.Matcher
}

func (s *seqscan) iter() {
}

func (s *seqscan) filter(m db.Matcher) {
	s.matchers = append(s.matchers, m)
}

func (s *seqscan) String() string {
	return fmt.Sprintf("seqscan{%v, %v}", s.modelID, s.matchers)
}

// dynamic join op
// nested loop join if less than 10 rows on one side
// hash join otherwise
type join struct {
	left, right relation
	leftBinding db.Binding
}

func (j *join) String() string {
	return fmt.Sprintf("join{%v, %v}", j.left, j.right)
}

func (j *join) iter() {
}

// type sort struct {
// }

type hashaggregate struct {
	r relation
	b db.Binding
	a Aggregation
}

func (a *hashaggregate) iter() {
}

func (a *hashaggregate) String() string {
	return fmt.Sprintf("hashaggregate{%v, %v}", a.r, a.a)
}

type limit struct {
	r     relation
	limit int
}

type setop int

const (
	union setop = iota
	subtraction
	intersection
)

type hashsetop struct {
	left, right relation
	op          setop
}

func (s *hashsetop) String() string {
	return fmt.Sprintf("hashsetop{%v, %v, %v}", s.left, s.right, s.op)
}

func (s *hashsetop) iter() {
}

func Plan(fm FindManyOperation) relation {
	s := planWhere(fm.ModelID, fm.Where)
	return s
}

func planWhere(modelID uuid.UUID, w Where) relation {

	// start with a scan, filtering on the table
	var s scan
	s = &seqscan{
		modelID: modelID,
	}

	// first we filter locally
	for _, fc := range w.FieldCriteria {
		s.filter(fc.Matcher())
	}
	var r relation
	r = s

	// then we join to-one related fields
	for _, rc := range w.RelationshipCriteria {
		rr := planRC(rc)
		r = &join{r, rr, rc.Binding}
	}

	// then we join to-many related fields
	for _, arc := range w.AggregateRelationshipCriteria {
		ar := planARC(arc)
		r = &join{r, ar, arc.RelationshipCriterion.Binding}
	}

	// then we union by or
	for _, orWhere := range w.Or {
		r = &hashsetop{r, planWhere(modelID, orWhere), union}
	}

	// then we subtract not
	// should this be conceptualized as an "antijoin" in the typical case
	for _, notWhere := range w.Not {
		r = &hashsetop{r, planWhere(modelID, notWhere), subtraction}
	}

	// then we intersect and
	for _, andWhere := range w.And {
		r = &hashsetop{r, planWhere(modelID, andWhere), intersection}
	}
	return r
}

func planRC(rc RelationshipCriterion) relation {
	d := rc.Binding.Dual()
	return planWhere(d.ModelID(), rc.Where)
}

func planARC(arc AggregateRelationshipCriterion) relation {
	b := arc.RelationshipCriterion.Binding

	r := planRC(arc.RelationshipCriterion)

	r = &hashaggregate{r, b, arc.Aggregation}
	return r
}
