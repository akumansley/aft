package db

import (
	"awans.org/aft/er/q"
	"awans.org/aft/internal/model"
	"github.com/google/uuid"
	"github.com/ompluscator/dynamic-struct"
	"reflect"
)

func getId(st interface{}) uuid.UUID {
	reader := dynamicstruct.NewReader(st)
	id := reader.GetField("Id").Interface().(uuid.UUID)
	return id
}

func newId(st *interface{}) {
	u := uuid.New()
	model.SystemAttrs["id"].SetField("id", u, *st)
}

func (op CreateOperation) Apply(db DB) interface{} {
	newId(&op.Struct)
	db.h.Insert(op.Struct)
	for _, no := range op.Nested {
		no.ApplyNested(db, op.Struct)
	}
	return op.Struct
}

// TODO hack -- remove this and rewrite with Relationship containing the name
func getBackref(db DB, rel model.Relationship) model.Relationship {
	m, _ := db.GetModel(rel.TargetModel)
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
	if fromRel.RelType == model.BelongsTo && (toRel.RelType == model.HasOne || toRel.RelType == model.HasMany) {
		// FK from
		setFK(from, toRel.TargetRel, getId(to))
	} else if toRel.RelType == model.BelongsTo && (fromRel.RelType == model.HasOne || fromRel.RelType == model.HasMany) {
		// FK to
		setFK(to, fromRel.TargetRel, getId(from))
	} else if toRel.RelType == model.HasManyAndBelongsToMany && fromRel.RelType == model.HasManyAndBelongsToMany {
		// Join table
		panic("Many to many relationships not implemented yet")
	} else {
		panic("Trying to connect invalid relationship")
	}
}

func (op NestedCreateOperation) ApplyNested(db DB, parent interface{}) {
	connect(db, parent, op.Relationship, op.Struct)
	newId(&op.Struct)
	db.h.Insert(op.Struct)
	db.h.Insert(parent)
}

func findOne(db DB, modelName string, uq UniqueQuery) interface{} {
	val, err := db.h.FindOne(modelName, q.Eq(uq.Key, uq.Val))
	if err != nil {
		panic("FindOne failed")
	}
	return val
}

func findOneById(db DB, modelName string, id uuid.UUID) interface{} {
	return findOne(db, modelName, UniqueQuery{Key: "Id", Val: id})
}

func (op NestedConnectOperation) ApplyNested(db DB, parent interface{}) {
	modelName := op.Relationship.TargetModel
	st := findOne(db, modelName, op.UniqueQuery)
	if st == nil {
		panic("No struct found")
	}
	connect(db, parent, op.Relationship, st)
	db.h.Insert(st)
	db.h.Insert(parent)
}
