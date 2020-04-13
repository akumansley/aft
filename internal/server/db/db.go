package db

import (
	"awans.org/aft/er"
	"awans.org/aft/er/q"
	"awans.org/aft/internal/model"
	"fmt"
	"reflect"
	"strings"
)

func New() DB {
	return DB{h: er.New()}
}

type DB struct {
	h *er.Hold
}

func (db DB) GetModel(modelName string) model.Model {
	modelName = strings.ToLower(modelName)
	val, err := db.h.FindOne("model", q.Eq("Name", modelName))
	if err != nil {
		panic(err)
	}
	m, ok := val.(model.Model)
	fmt.Printf("val is %v\n", val)
	if !ok {
		panic("Not a model")
	}
	return m
}

func (db DB) MakeStruct(modelName string) interface{} {
	modelName = strings.ToLower(modelName)
	if modelName == "model" {
		return model.Model{}
	} else {
		m := db.GetModel(modelName)
		st := model.StructForModel(m).New()
		field := reflect.ValueOf(st).Elem().FieldByName("Type")
		field.SetString(modelName)
		return st
	}
}
