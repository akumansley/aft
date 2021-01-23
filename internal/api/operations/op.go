package operations

import (
	"encoding/gob"

	"awans.org/aft/internal/db"
)

type FindArgs struct {
	Where   Where     `json:",omitempty"`
	Include Include   `json:",omitempty"`
	Select  Select    `json:",omitempty"`
	Case    Case      `json:",omitempty"`
	Order   []db.Sort `json:",omitempty"`
	Take    int
	Skip    int
}

func init() {
	gob.Register(FindArgs{})
	gob.Register(FindOneOperation{})
	gob.Register(FindManyOperation{})
	gob.Register(CreateOperation{})
	gob.Register(UpdateOperation{})
	gob.Register(UpdateManyOperation{})
	gob.Register(UpsertOperation{})
	gob.Register(DeleteOperation{})
	gob.Register(DeleteManyOperation{})
	gob.Register(CountOperation{})
	gob.Register(NestedCreateOperation{})
	gob.Register(NestedConnectOperation{})
	gob.Register(NestedDisconnectOperation{})
	gob.Register(NestedSetOperation{})
	gob.Register(NestedUpdateOperation{})
	gob.Register(NestedUpdateManyOperation{})
	gob.Register(NestedDeleteOperation{})
	gob.Register(NestedDeleteManyOperation{})
	gob.Register(NestedUpsertOperation{})
}

type FindOneOperation struct {
	ModelID  db.ID    `json:",omitempty"`
	FindArgs FindArgs `json:",omitempty"`
}

type FindManyOperation struct {
	ModelID  db.ID    `json:",omitempty"`
	FindArgs FindArgs `json:",omitempty"`
}

type CreateOperation struct {
	ModelID  db.ID                  `json:",omitempty"`
	Data     map[string]interface{} `json:",omitempty"`
	FindArgs FindArgs               `json:",omitempty"`
	Nested   []NestedOperation      `json:",omitempty"`
}

type UpdateOperation struct {
	ModelID  db.ID                  `json:",omitempty"`
	FindArgs FindArgs               `json:",omitempty"`
	Data     map[string]interface{} `json:",omitempty"`
	Nested   []NestedOperation      `json:",omitempty"`
}

type UpsertOperation struct {
	ModelID      db.ID                  `json:",omitempty"`
	FindArgs     FindArgs               `json:",omitempty"`
	Create       map[string]interface{} `json:",omitempty"`
	NestedCreate []NestedOperation      `json:",omitempty"`
	Update       map[string]interface{} `json:",omitempty"`
	NestedUpdate []NestedOperation      `json:",omitempty"`
}

type DeleteOperation struct {
	ModelID  db.ID             `json:",omitempty"`
	FindArgs FindArgs          `json:",omitempty"`
	Nested   []NestedOperation `json:",omitempty"`
}

type UpdateManyOperation struct {
	ModelID db.ID                  `json:",omitempty"`
	Where   Where                  `json:",omitempty"`
	Data    map[string]interface{} `json:",omitempty"`
	Nested  []NestedOperation      `json:",omitempty"`
}

type DeleteManyOperation struct {
	ModelID db.ID             `json:",omitempty"`
	Where   Where             `json:",omitempty"`
	Nested  []NestedOperation `json:",omitempty"`
}

type CountOperation struct {
	ModelID db.ID `json:",omitempty"`
	Where   Where `json:",omitempty"`
}

//Nested operations
type NestedOperation interface {
	ApplyNested(tx db.RWTx, parent db.ModelRef, parents []*db.QueryResult) error
}

type NestedCreateOperation struct {
	Relationship db.Relationship        `json:",omitempty"`
	Model        db.Model               `json:",omitempty"`
	Data         map[string]interface{} `json:",omitempty"`
	Nested       []NestedOperation      `json:",omitempty"`
}

type NestedConnectOperation struct {
	Relationship db.Relationship `json:",omitempty"`
	Where        Where           `json:",omitempty"`
}

type NestedDisconnectOperation struct {
	Relationship db.Relationship `json:",omitempty"`
	Where        Where           `json:",omitempty"`
}

type NestedSetOperation struct {
	Relationship db.Relationship `json:",omitempty"`
	Where        Where           `json:",omitempty"`
}

type NestedUpdateOperation struct {
	Relationship db.Relationship        `json:",omitempty"`
	Where        Where                  `json:",omitempty"`
	Data         map[string]interface{} `json:",omitempty"`
	Nested       []NestedOperation      `json:",omitempty"`
}

type NestedDeleteOperation struct {
	Relationship db.Relationship   `json:",omitempty"`
	Where        Where             `json:",omitempty"`
	Nested       []NestedOperation `json:",omitempty"`
}

type NestedUpdateManyOperation struct {
	Relationship db.Relationship        `json:",omitempty"`
	Where        Where                  `json:",omitempty"`
	Data         map[string]interface{} `json:",omitempty"`
	Nested       []NestedOperation      `json:",omitempty"`
}

type NestedDeleteManyOperation struct {
	Relationship db.Relationship   `json:",omitempty"`
	Where        Where             `json:",omitempty"`
	Nested       []NestedOperation `json:",omitempty"`
}

type NestedUpsertOperation struct {
	Relationship db.Relationship        `json:",omitempty"`
	Where        Where                  `json:",omitempty"`
	Create       map[string]interface{} `json:",omitempty"`
	NestedCreate []NestedOperation      `json:",omitempty"`
	Update       map[string]interface{} `json:",omitempty"`
	NestedUpdate []NestedOperation      `json:",omitempty"`
}
