package db

import (
	"github.com/google/uuid"
)

type EnumValue struct {
	ID   ID
	Name string
}

type RuntimeEnumValue struct {
	EnumValue
}
type StorageEnumValue struct {
	EnumValue
}
type FunctionSignatureEnumValue struct {
	EnumValue
}

func RecordToEnumValue(r Record, field string, tx Tx) (EnumValue, error) {
	v, err := r.Get(field)
	if err != nil {
		return EnumValue{}, err
	}
	id := v.(uuid.UUID)
	rec, err := tx.FindOne(EnumValueModel.ID, EqID(ID(id)))
	if err != nil {
		return EnumValue{}, err
	}
	name, err := rec.Get("name")
	if err != nil {
		return EnumValue{}, err
	}
	return EnumValue{ID: ID(id), Name: name.(string)}, nil
}

var enumMap map[ID]EnumValue = map[ID]EnumValue{
	Native.ID:        Native.EnumValue,
	Starlark.ID:      Starlark.EnumValue,
	FromJSON.ID:      FromJSON.EnumValue,
	RPC.ID:           RPC.EnumValue,
	BoolStorage.ID:   BoolStorage.EnumValue,
	IntStorage.ID:    IntStorage.EnumValue,
	StringStorage.ID: StringStorage.EnumValue,
	FloatStorage.ID:  FloatStorage.EnumValue,
	UUIDStorage.ID:   UUIDStorage.EnumValue,
}
