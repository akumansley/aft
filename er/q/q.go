package q

import (
	"reflect"
)

type Matcher interface {
	Match(interface{}) (bool, error)
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

func getFieldIf(field string, st interface{}) interface{} {
	k := reflect.ValueOf(st).Kind()
	switch k {
	case reflect.Struct:
		return reflect.ValueOf(st).FieldByName(field).Interface()
	case reflect.Interface:
		return reflect.ValueOf(st).Elem().FieldByName(field).Interface()
	case reflect.Ptr:
		return reflect.ValueOf(st).Elem().FieldByName(field).Interface()

	}
	return nil
}

// could be faster probably
func (fm FieldMatcher) Match(st interface{}) (bool, error) {
	candidate := getFieldIf(fm.field, st)
	comparison := fm.val
	return candidate == comparison, nil
}

func Eq(field string, val interface{}) Matcher {
	return FieldMatcher{field: field, val: val, op: eq}
}
