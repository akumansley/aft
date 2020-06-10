package api

import (
	"awans.org/aft/internal/db"
	"fmt"
	"github.com/google/uuid"
)

type iterator interface {
	Next() bool
	Value() interface{}
	Err() error
}

type reciter struct {
	recs  []db.Record
	ix    int
	value db.Record
	err   error
}

func (i *reciter) Value() interface{} {
	return i.value
}
func (i *reciter) Err() error {
	return i.err
}

func (i *reciter) Next() bool {
	if i.ix < len(i.recs) {
		i.ix++
		i.value = i.recs[i.ix-1]
		return true
	}
	return false
}

type relation interface {
	iter(db.Tx) iterator
	String() string
}

type scan struct {
	modelID  uuid.UUID
	ID       int
	matchers []db.Matcher
}

func (s *scan) iter(tx db.Tx) iterator {
	recs := tx.FindMany(s.modelID, db.And(s.matchers...))
	return &reciter{recs: recs, ix: 0}
}

func (s *scan) filter(m db.Matcher) {
	s.matchers = append(s.matchers, m)
}

func (s *scan) String() string {
	return fmt.Sprintf("scan{%v, %v}", s.modelID, s.matchers)
}

type frame struct {
	recs []db.Record
}

type innerjoin1 struct {
	left, right relation
	lsid, rsid  int
	leftBinding db.Binding
}

func (j *innerjoin) String() string {
	return fmt.Sprintf("join{%v, %v}", j.left, j.right)
}

func (j *innerjoin) iter(tx db.Tx) iterator {
	relType := j.leftBinding.RelType()
	switch relType {
	case db.HasOne:
		// FK on the right side
		// start by iterating on the right side
	case db.BelongsTo:
		// FK on the left side
	case db.HasMany:
		// FK on the right side
	case db.HasManyAndBelongsToMany:
		panic("Not implemented")
	}

	leftIter := j.left.iter(tx)

	for leftIter.Next() {
		val := leftIter.Value()
	}
}

// type sort struct {
// }

type hashaggregate struct {
	r relation
	b db.Binding
	a Aggregation
}

func (a *hashaggregate) iter(tx db.Tx) iterator {
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
	lsid, rsid  int
	op          setop
}

func (s *hashsetop) String() string {
	return fmt.Sprintf("hashsetop{%v, %v, %v}", s.left, s.right, s.op)
}

func (s *hashsetop) iter(tx db.Tx) iterator {
}

type framemaker struct {
	capacity int
}

func (fm framemaker) makeFrame() *frame {
	return &frame{
		recs: make([]db.Record, fm.capacity),
	}
}

type ExecutionPlan struct {
	r        relation
	numScans int
}

func (e ExecutionPlan) Execute(tx db.Tx) {
	fm := &framemaker{e.numScans}
	// TODO
}

func Plan(fm FindManyOperation) ExecutionPlan {
	jt := &jointree{}
	_, r := planWhere(fm.ModelID, fm.Where, jt)
	return ExecutionPlan{r, *c}
}

type jointree struct {
	nextSid  int
	children map[int]map[db.Binding]int
}

func (j *jointree) nextID() int {
	sid := j.nextSid
	j.nextSid++
	return sid
}

func (j *jointree) join(parent, child int, b db.Binding) {
	jm, ok := j.children[parent]
	if !ok {
		j.children[parent] = map[db.Binding]int{
			b: child,
		}
	} else {
		jm[b] = child
	}
}

func planWhere(modelID uuid.UUID, w Where, jt *jointree) (int, relation) {

	// start with a scan, filtering on the table
	// the "scan id" (sid) functions as an alias of sorts
	// so we can differentiate between different scans of a self-join
	// we increment it for each new scan
	sid := jt.nextID()
	s := &scan{
		modelID: modelID,
		ID:      sid,
	}

	// first we filter locally
	for _, fc := range w.FieldCriteria {
		s.filter(fc.Matcher())
	}
	var r relation
	r = s

	// then we join to-one related fields
	for _, rc := range w.RelationshipCriteria {
		rcSid, rr := planRC(rc, jt)
		r = &innerjoin{r, rr, sid, rcSid, rc.Binding}
		jt.join(sid, rcSid, rc.Binding)
	}

	// then we join to-many related fields
	for _, arc := range w.AggregateRelationshipCriteria {
		arSid, ar := planARC(arc, jt)
		r = &innerjoin{r, ar, sid, arSid, arc.RelationshipCriterion.Binding}
		jt.join(sid, arSid, arc.RelationshipCriterion.Binding)
	}

	// then we union by or
	for _, orWhere := range w.Or {
		orSid, orR := planWhere(modelID, orWhere, jt)
		r = &hashsetop{r, orR, sid, orSid, union}
	}

	// then we subtract not
	// should this be conceptualized as an "antijoin" in the typical case
	for _, notWhere := range w.Not {
		notSid, notR := planWhere(modelID, notWhere, jt)
		r = &hashsetop{r, notR, sid, notSid, subtraction}
	}

	// then we intersect and
	for _, andWhere := range w.And {
		andSid, andR := planWhere(modelID, andWhere, jt)
		r = &hashsetop{r, andR, sid, andSid, intersection}
	}
	return sid, r
}

func planRC(rc RelationshipCriterion, jt *jointree) (int, relation) {
	d := rc.Binding.Dual()
	return planWhere(d.ModelID(), rc.Where, jt)
}

func planARC(arc AggregateRelationshipCriterion, jt *jointree) (int, relation) {
	b := arc.RelationshipCriterion.Binding

	rSid, r := planRC(arc.RelationshipCriterion, jt)

	r = &hashaggregate{r, b, arc.Aggregation}
	return rSid, r
}
