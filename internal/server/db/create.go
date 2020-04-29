package db

import (
	"awans.org/aft/er/q"
	"awans.org/aft/internal/model"
	"fmt"
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

func (op CreateOperation) Apply(db DB) (interface{}, error) {
	newId(&op.Struct)
	db.h.Insert(op.Struct)
	for _, no := range op.Nested {
		no.ApplyNested(db, op.Struct)
	}
	return op.Struct, nil
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
		fmt.Printf("fromRel %v toRel %v\n", fromRel, toRel)
		panic("Trying to connect invalid relationship")
	}
}

func (op NestedCreateOperation) ApplyNested(db DB, parent interface{}) (err error) {
	connect(db, parent, op.Relationship, op.Struct)
	newId(&op.Struct)
	db.h.Insert(op.Struct)
	db.h.Insert(parent)
	return nil
}

func findOne(db DB, modelName string, uq UniqueQuery) (st interface{}, err error) {
	st, err = db.h.FindOne(modelName, q.Eq(uq.Key, uq.Val))
	return
}

func findOneById(db DB, modelName string, id uuid.UUID) (st interface{}, err error) {
	return findOne(db, modelName, UniqueQuery{Key: "Id", Val: id})
}

func (op NestedConnectOperation) ApplyNested(db DB, parent interface{}) (err error) {
	modelName := op.Relationship.TargetModel
	st, err := findOne(db, modelName, op.UniqueQuery)
	if err != nil {
		return
	}
	connect(db, parent, op.Relationship, st)
	db.h.Insert(st)
	db.h.Insert(parent)
	return
}
