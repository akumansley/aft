package db

import (
	"awans.org/aft/internal/model"
	"github.com/google/uuid"
	"github.com/ompluscator/dynamic-struct"
	"github.com/timshannon/bolthold"
	"reflect"
)

func getId(st interface{}) uuid.UUID {
	reader := dynamicstruct.NewReader(st)
	id := reader.GetField("Id").Interface().(uuid.UUID)
	return id
}

func insert(db DB, st interface{}) {
	db.db.Insert(getId(st), st)
}

func (op CreateOperation) Apply(db DB) {
	insert(db, op.Struct)
	for _, no := range op.Nested {
		no.ApplyNested(db, op.Struct)
	}
}

func getBackref(db DB, rel model.Relationship) model.Relationship {
	m := db.GetModel(rel.TargetModel)
	return m.Relationships[rel.TargetRel]
}

func setFK(st interface{}, key string, id uuid.UUID) {
	fieldName := model.JsonKeyToRelFieldName(key)
	field := reflect.ValueOf(st).Elem().FieldByName(fieldName)
	v := reflect.ValueOf(id)
	field.Set(v)
}

func connect(db DB, from interface{}, fromRel model.Relationship, to interface{}) {
	toRel := getBackref(db, fromRel)
	if fromRel.Type == model.BelongsTo && (toRel.Type == model.HasOne || toRel.Type == model.HasMany) {
		// FK from
		setFK(from, toRel.TargetRel, getId(to))
	} else if toRel.Type == model.BelongsTo && (fromRel.Type == model.HasOne || fromRel.Type == model.HasMany) {
		// FK to
		setFK(to, fromRel.TargetRel, getId(from))
	} else if toRel.Type == model.HasManyAndBelongsToMany && fromRel.Type == model.HasManyAndBelongsToMany {
		// Join table
		panic("Many to many relationships not implemented yet")
	} else {
		panic("Trying to connect invalid relationship")
	}
}

func (op NestedCreateOperation) ApplyNested(db DB, parent interface{}) {
	connect(db, parent, op.Relationship, op.Struct)
	insert(db, op.Struct)
	insert(db, parent)
}

func findOne(db DB, st *interface{}, uq UniqueQuery) {
	err := db.db.FindOne(st, bolthold.Where(uq.Key).Eq(uq.Val))
	if err != nil {
		panic("FindOne failed")
	}
}

func (op NestedConnectOperation) ApplyNested(db DB, parent interface{}) {
	st := db.MakeStruct(op.Relationship.TargetModel)
	findOne(db, &st, op.UniqueQuery)
	connect(db, parent, op.Relationship, st)
	insert(db, st)
	insert(db, parent)
}
