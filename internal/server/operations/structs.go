package operations

import (
	"awans.org/aft/internal/model"
)

type Operation interface {
}
type NestedOperation interface {
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
	Val string
}

type NestedConnectOperation struct {
	Relationship model.Relationship
	UniqueQuery  UniqueQuery
}
