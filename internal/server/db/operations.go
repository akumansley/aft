package db

import (
	"awans.org/aft/internal/model"
)

type Operation interface {
	Apply(DB)
}

type NestedOperation interface {
	ApplyNested(DB, interface{})
}

type CreateOperation struct {
	Struct interface{}
	Nested []NestedOperation
}

type NestedCreateOperation struct {
	Relationship model.Relationship
	Struct       interface{}
	Nested       []NestedOperation
}

type UniqueQuery struct {
	Key string
	Val interface{}
}

type NestedConnectOperation struct {
	Relationship model.Relationship
	UniqueQuery  UniqueQuery
}
