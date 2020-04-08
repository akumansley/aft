package operations

type Operation interface {
}

type CreateOperation struct {
	Struct interface{}
	Nested []NestedCreateOperation
}

type NestedCreateOperation struct {
	Parent       interface{}
	Relationship string
	Struct       interface{}
	Nested       []NestedCreateOperation
}

type NestedConnectOperation struct {
	Parent       interface{}
	Relationship string
	Struct       interface{}
	Nested       []NestedCreateOperation
}
