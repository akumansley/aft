package db

import (
	"github.com/google/uuid"
)

type Aggregation int

const (
	Every Aggregation = iota
	Some
	None
)

type ModelRef struct {
	modelID uuid.UUID
	aliasID uuid.UUID
}

type RefMatcher struct {
	ref     ModelRef
	matcher Matcher
}

func (m ModelRef) Eq(key string, val interface{}) RefMatcher {
	return RefMatcher{m, Eq(key, val)}
}

func (m ModelRef) Match(match Matcher) RefMatcher {
	return RefMatcher{m, match}
}

func Ref(modelID uuid.UUID) ModelRef {
	return ModelRef{modelID, uuid.New()}
}

type relation interface {
	hasAlias(ModelRef) bool
	plan(*Query) PlanNode
}

type table struct {
	ref ModelRef
}

func (t table) hasAlias(m ModelRef) bool {
	return m == t.ref
}

type OnPredicate struct {
	left, right ModelRef
	b           Binding
}

func On(left, right ModelRef, b Binding) OnPredicate {
	return OnPredicate{left, right, b}
}

type joinone struct {
	left, right relation
	predicate   OnPredicate
}

func (j joinone) hasAlias(m ModelRef) bool {
	return j.left.hasAlias(m) || j.right.hasAlias(m)
}

type setoperation int

const (
	union setoperation = iota
	intersection
	exclusion
)

type setop struct {
	left, right relation

	// must be two refs to the same model
	// and have disjoint aliases... or things will get messed up
	leftRef, rightRef ModelRef
	op                setoperation
}

func (s setop) hasAlias(m ModelRef) bool {
	return s.left.hasAlias(m) || s.right.hasAlias(m)
}

type joinmany struct {
	left, right relation
	predicate   OnPredicate
	agg         Aggregation
}

func (j joinmany) hasAlias(m ModelRef) bool {
	return j.left.hasAlias(m) || j.right.hasAlias(m)
}

type Query struct {
	r relation
	// map by alias id
	predicateMap map[uuid.UUID][]Matcher
	aset         map[uuid.UUID]bool
}

func (q *Query) hasAlias(a ModelRef) bool {
	_, ok := q.aset[a.aliasID]
	return ok
}

func Select(a ModelRef) *Query {
	q := Query{r: table{a}}
	q.predicateMap = map[uuid.UUID][]Matcher{}
	q.aset = map[uuid.UUID]bool{a.aliasID: true}
	return &q
}

func (q *Query) Where(m RefMatcher) {
	if !q.hasAlias(m.ref) {
		panic("bad query")
	}
	key := m.ref.aliasID
	filters, ok := q.predicateMap[key]
	if !ok {
		q.predicateMap[key] = []Matcher{m.matcher}
	} else {
		q.predicateMap[key] = append(filters, m.matcher)
	}
}

func (q *Query) JoinOne(right *Query, o OnPredicate) *Query {
	// also assert the OnPredicate has the right sort of binding
	if !q.hasAlias(o.left) {
		panic("bad query")
	}

	if !right.hasAlias(o.right) {
		panic("bad query")
	}

	q.r = &joinone{q.r, right.r, o}
	for k := range right.aset {
		_, ok := q.aset[k]
		if ok {
			panic("non-disjoint aliases in join")
		}
		q.aset[k] = true
	}

	return q
}

func (q *Query) JoinMany(right *Query, o OnPredicate, a Aggregation) *Query {
	// also assert the OnPredicate has the right sort of binding
	if !q.hasAlias(o.left) {
		panic("bad query")
	}

	if !right.hasAlias(o.right) {
		panic("bad query")
	}

	q.r = &joinmany{q.r, right.r, o, a}

	for k := range right.aset {
		_, ok := q.aset[k]
		if ok {
			panic("non-disjoint aliases in join")
		}
		q.aset[k] = true
	}

	return q
}

func (q *Query) hashsetop(other *Query, left, right ModelRef, op setoperation) *Query {
	if left.modelID != right.modelID {
		panic("bad query")
	}
	if left.aliasID == right.aliasID {
		panic("bad query")
	}

	q.r = &setop{q.r, other.r, left, right, op}

	for k := range other.aset {
		_, ok := q.aset[k]
		if ok {
			panic("non-disjoint aliases in hashsetop")
		}
		q.aset[k] = true
	}

	return q
}

func (q *Query) Union(other *Query, left, right ModelRef) *Query {
	return q.hashsetop(other, left, right, union)
}

func (q *Query) Intersect(other *Query, left, right ModelRef) *Query {
	return q.hashsetop(other, left, right, intersection)
}

func (q *Query) Except(other *Query, left, right ModelRef) *Query {
	return q.hashsetop(other, left, right, exclusion)
}

// m1 := db.Ref(modelID)
// m2 := db.Ref(modelID2)
// q := Select(m1).Where(m1.Eq(k, v)).Join(m1, m2, db.Binding)
// q4 := q.Union(uq)
// q5 := q.Insersection(uq)
