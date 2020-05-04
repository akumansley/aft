package db

import (
	"awans.org/aft/internal/model"
	"github.com/google/uuid"
	"github.com/ompluscator/dynamic-struct"
)

func (i Include) Resolve(tx Tx, st interface{}) interface{} {
	for _, inc := range i.Includes {
		tx.Resolve(st, inc)
	}
	return st
}

func getFK(st interface{}, key string) uuid.UUID {
	fieldName := model.JsonKeyToRelFieldName(key)
	reader := dynamicstruct.NewReader(st)
	id := reader.GetField(fieldName).Interface().(uuid.UUID)
	return id
}
