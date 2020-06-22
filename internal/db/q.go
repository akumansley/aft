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
	candidate := st.MustGet(fm.field)
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

type FKFieldMatcher struct {
	field string
	val   ID
	op    op
}

func (fm FKFieldMatcher) Match(st Record) (bool, error) {
	candidate := st.MustGetFK(fm.field)
	comparison := fm.val
	return candidate == comparison, nil
}

func EqFK(field string, val ID) Matcher {
	return FKFieldMatcher{field: field, val: val, op: eq}
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
