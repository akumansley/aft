package db

import (
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

func (op CreateOperation) Apply(tx RWTx) (interface{}, error) {
	newId(&op.Struct)
	tx.Insert(op.Struct)
	for _, no := range op.Nested {
		no.ApplyNested(tx, op.Struct)
	}
	return op.Struct, nil
}
func setFK(st interface{}, key string, id uuid.UUID) {
	fieldName := model.JsonKeyToRelFieldName(key)
	field := reflect.ValueOf(st).Elem().FieldByName(fieldName)
	v := reflect.ValueOf(id)
	field.Set(v)
}

func (op NestedCreateOperation) ApplyNested(tx RWTx, parent interface{}) (err error) {
	newId(&op.Struct)
	tx.Connect(parent, op.Struct, op.Relationship)
	return nil
}

func findOneById(tx Tx, modelName string, id uuid.UUID) (st interface{}, err error) {
	return tx.FindOne(modelName, UniqueQuery{Key: "Id", Val: id})
}

func (op NestedConnectOperation) ApplyNested(tx RWTx, parent interface{}) (err error) {
	modelName := op.Relationship.TargetModel
	st, err := tx.FindOne(modelName, op.UniqueQuery)
	if err != nil {
		return
	}
	tx.Connect(parent, st, op.Relationship)
	return
}
