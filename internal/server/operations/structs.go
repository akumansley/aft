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

type NestedConnectOperation struct {
	Relationship string
	UniqueQuery  interface{}
}
