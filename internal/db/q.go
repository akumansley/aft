package db

import (
	"awans.org/aft/internal/model"
	"github.com/google/uuid"
)

type Matcher interface {
	Match(model.Record) (bool, error)
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
func (fm FieldMatcher) Match(st model.Record) (bool, error) {
	candidate := st.Get(fm.field)
	comparison := fm.val
	return candidate == comparison, nil
}

func Eq(field string, val interface{}) Matcher {
	return FieldMatcher{field: field, val: val, op: eq}
}

type FKFieldMatcher struct {
	field string
	val   uuid.UUID
	op    op
}

func (fm FKFieldMatcher) Match(st model.Record) (bool, error) {
	candidate := st.GetFK(fm.field)
	comparison := fm.val
	return candidate == comparison, nil
}

func EqFK(field string, val uuid.UUID) Matcher {
	return FKFieldMatcher{field: field, val: val, op: eq}
}

type AndMatcher struct {
	inner []Matcher
}

func (am AndMatcher) Match(st model.Record) (bool, error) {
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
