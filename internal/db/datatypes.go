package db

import (
	"awans.org/aft/internal/datatypes"
	"fmt"
	"github.com/google/uuid"
)

type Datatype struct {
	ID        ID
	Name      string
	StoredAs  Storage
	Validator Code
}

type Storage int64

const (
	BoolStorage Storage = iota
	IntStorage
	StringStorage
	FloatStorage
	UUIDStorage
)

var storageMap map[Storage]interface{} = map[Storage]interface{}{
	BoolStorage:   false,
	IntStorage:    int64(0),
	StringStorage: "",
	FloatStorage:  0.0,
	UUIDStorage:   uuid.UUID{},
}

func (d Datatype) typeCheck(out interface{}) (interface{}, error) {
	switch d.StoredAs {
	case BoolStorage:
		return datatypes.BoolFromJSON(out)
	case IntStorage:
		return datatypes.IntFromJSON(out)
	case StringStorage:
		return datatypes.StringFromJSON(out)
	case FloatStorage:
		return datatypes.FloatFromJSON(out)
	case UUIDStorage:
		return datatypes.UUIDFromJSON(out)
	}
	return nil, fmt.Errorf("Unrecognized storage type for datatype")
}

var datatypeMap map[ID]Datatype = map[ID]Datatype{
	Bool.ID:   Bool,
	Int.ID:    Int,
	Enum.ID:   Enum,
	String.ID: String,
	Text.ID:   Text,
	UUID.ID:   UUID,
	Float.ID:  Float,
}

func (d Datatype) FromJSON(arg interface{}) (interface{}, error) {
	c := d.Validator
	out, err := c.executor.Invoke(c, arg)
	if err != nil {
		return nil, err
	}
	return d.typeCheck(out)
}
