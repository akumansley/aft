package operations

import (
	"awans.org/aft/internal/db"
)

type FindArgs struct {
	Where   Where
	Include Include
	// Add Select
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
	Record   db.Record
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
	Create       db.Record
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
	ApplyNested(tx db.RWTx, root db.ModelRef, parent db.ModelRef, parents []*db.QueryResult, clauses []db.QueryClause) error
}

type NestedCreateOperation struct {
	Relationship db.Relationship
	Record       db.Record
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
	Create       db.Record
	NestedCreate []NestedOperation
	Update       map[string]interface{}
	NestedUpdate []NestedOperation
}
