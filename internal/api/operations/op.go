package operations

import (
	"awans.org/aft/internal/db"
)

type FindArgs struct {
	Where   Where
	Include Include
	Select  Select
}

type FindOneOperation struct {
	ModelID  db.ID
	FindArgs FindArgs
}

type FindManyOperation struct {
	ModelID  db.ID
	FindArgs FindArgs
}

type CreateOperation struct {
	ModelID  db.ID
	Data     map[string]interface{}
	FindArgs FindArgs
	Nested   []NestedOperation
}

type UpdateOperation struct {
	ModelID  db.ID
	FindArgs FindArgs
	Data     map[string]interface{}
	Nested   []NestedOperation
}

type UpsertOperation struct {
	ModelID      db.ID
	FindArgs     FindArgs
	Create       map[string]interface{}
	NestedCreate []NestedOperation
	Update       map[string]interface{}
	NestedUpdate []NestedOperation
}

type DeleteOperation struct {
	ModelID  db.ID
	FindArgs FindArgs
	Nested   []NestedOperation
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
	ApplyNested(tx db.RWTx, parent db.ModelRef, parents []*db.QueryResult) error
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

type NestedDisconnectOperation struct {
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
