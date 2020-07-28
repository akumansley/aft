package operations

import (
	"awans.org/aft/internal/db"
)

type FindArgs struct {
	Where   Where
	Include Include
	// Add Select
}

type FindManyOperation struct {
	ModelID db.ID
	FindManyArgs FindManyArgs
}

type FindManyArgs struct {
	Where   Where
	Include Include
	// Add Select
}

type CreateOperation struct {
	Record  db.Record
	Include Include
	// Add Select
	Nested  []NestedOperation
}

type UpdateOperation struct {
	ModelID db.ID
	Where   Where
	Data    map[string]interface{}
	Include Include
	//Add Select
	Nested []NestedOperation
}

type UpsertOperation struct {
	ModelID      db.ID
	Where        Where
	Create       db.Record
	NestedCreate []NestedOperation
	Update       map[string]interface{}
	NestedUpdate []NestedOperation
	Include      Include
	// Add Select
}

type DeleteOperation struct {
	Where   Where
	ModelID db.ID
	Include Include
	// Add Select
	Nested  []NestedOperation
}

type UpdateManyOperation struct {
	ModelID db.ID
	Where   Where
	Data    map[string]interface{}
	Nested  []NestedOperation
}

type DeleteManyOperation struct {
	ModelID db.ID
	Where   Where
	Nested  []NestedOperation
}

type CountOperation struct {
	ModelID db.ID
	Where   Where
}

//Nested operations
type NestedOperation interface {
	ApplyNested(db.RWTx) error
}

type NestedCreateOperation struct {
	Relationship db.Relationship
	Data         map[string]interface{}
	Nested       []NestedOperation
}

type NestedConnectOperation struct {
	Relationship db.Relationship
	Where        Where
}

type NestedUpdateOperation struct {
	Where        Where
	Relationship db.Relationship
	Data         map[string]interface{}
	Nested       []NestedOperation
}

type NestedDeleteOperation struct {
	Where        Where
	Relationship db.Relationship
	Nested       []NestedOperation
}

type NestedUpdateManyOperation struct {
	Where        Where
	Relationship db.Relationship
	Data         map[string]interface{}
	Nested       []NestedOperation
}

type NestedDeleteManyOperation struct {
	Where        Where
	Relationship db.Relationship
	Nested       []NestedOperation
}

type NestedUpsertOperation struct {
	Relationship db.Relationship
	Where        Where
	Create       map[string]interface{}
	NestedCreate []NestedOperation
	Update       map[string]interface{}
	NestedUpdate []NestedOperation
}
