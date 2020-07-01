package db

import (
	"github.com/google/uuid"
)

type Matcher interface {
	Match(Record) (bool, error)
}

type op int

const (
	eq op = iota
	neq
	gt // not implemented
	lt // not implemented
)

type FieldMatcher struct {
	field string
	val   interface{}
	op    op
}

// could be faster probably
func (fm FieldMatcher) Match(st Record) (bool, error) {
	candidate := st.mustGet(fm.field)
	comparison := fm.val
	return candidate == comparison, nil
}

func EqID(val ID) Matcher {
	u := uuid.UUID(val)
	return FieldMatcher{field: "id", val: u, op: eq}
}

func Eq(field string, val interface{}) Matcher {
	return FieldMatcher{field: field, val: val, op: eq}
}

type void struct{}

type idSetMatcher struct {
	ids map[uuid.UUID]void
}

func (im idSetMatcher) Match(r Record) (bool, error) {
	id := r.mustGet("id")
	_, ok := im.ids[id.(uuid.UUID)]
	return ok, nil
}

func IDIn(ids []ID) Matcher {
	hash := make(map[uuid.UUID]void)
	for _, id := range ids {
		u := uuid.UUID(id)
		hash[u] = void{}
	}

	return idSetMatcher{ids: hash}
}

type AndMatcher struct {
	inner []Matcher
}

func (am AndMatcher) Match(st Record) (bool, error) {
	for _, m := range am.inner {
		match, err := m.Match(st)
		if err != nil {
			return false, err
		}
		if !match {
			return false, nil
		}
	}
	return true, nil
}

func And(matchers ...Matcher) Matcher {
	return AndMatcher{inner: matchers}
}
