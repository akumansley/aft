package operations

import (
	"awans.org/aft/internal/db"
)

type FindOneOperation struct {
	ModelID db.ID
	Where   Where
	Include Include
	// Add Select
}

type CreateOperation struct {
	Record  db.Record
	Include Include
	Nested  []NestedOperation
}

type UpdateOperation struct {
	Old     db.Record
	New     db.Record
	Include Include
	Nested  []NestedOperation
}

type FindManyOperation struct {
	ModelID db.ID
	Where   Where
	Include Include
	// Add Select
}

type UpdateManyOperation struct {
	Old []db.Record
	New []db.Record
}

<<<<<<< HEAD
=======
type CountOperation struct {
	ModelID db.ID
	Where   Where
}

>>>>>>> 6b57a71... Refactor API to use common code
//Nested operations
type NestedOperation interface {
	ApplyNested(db.RWTx, db.Record) error
}

type NestedCreateOperation struct {
	Relationship db.Relationship
	Record       db.Record
	Nested       []NestedOperation
}

type NestedConnectOperation struct {
	Relationship db.Relationship
	UniqueQuery  UniqueQuery
}

type NestedUpdateOperation struct {
	Old    db.Record
	New    db.Record
	Nested []NestedOperation
}

type NestedUpdateManyOperation struct {
	Old []db.Record
	New []db.Record
}

type NestedFindManyOperation struct {
	Where   Where
	Include Include
	// Add Select
}
