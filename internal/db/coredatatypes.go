package db

import (
	"fmt"
	"strconv"

	"github.com/google/uuid"
)

var (
	ErrValue     = fmt.Errorf("invalid value for type")
	ErrNotStored = fmt.Errorf("value not stored")
)

func BoolFromJSON(args []interface{}) (interface{}, error) {
	value := args[0]
	b, ok := value.(bool)
	if !ok {
		return nil, fmt.Errorf("%w: expected bool got %T", ErrValue, value)
	}
	return b, nil
}

var boolValidator = MakeNativeFunction(
	MakeID("8e806967-c462-47af-8756-48674537a909"),
	"bool",
	1,
	BoolFromJSON)

func IntFromJSON(args []interface{}) (interface{}, error) {
	value := args[0]
	switch value.(type) {
	case int64:
		return value, nil
	case int:
		return int64(value.(int)), nil
	case float64:
		return int64(value.(float64)), nil
	}
	return nil, fmt.Errorf("%w: expected int got %T", ErrValue, value)

}

var intValidator = MakeNativeFunction(
	MakeID("a1cf1c16-040d-482c-92ae-92d59dbad46c"),
	"int",
	1,
	IntFromJSON)

func StringFromJSON(args []interface{}) (interface{}, error) {
	value := args[0]
	s, ok := value.(string)
	if !ok {
		return nil, fmt.Errorf("%w: expected string got %T", ErrValue, value)
	}
	return s, nil
}

var stringValidator = MakeNativeFunction(
	MakeID("aaeccd14-e69f-4561-91ef-5a8a75b0b498"),
	"string",
	1,
	StringFromJSON)

func UUIDFromJSON(args []interface{}) (interface{}, error) {
	value := args[0]
	var u uuid.UUID
	var err error
	switch value.(type) {
	case string:
		u, err = uuid.Parse(value.(string))
		if err != nil {
			return nil, fmt.Errorf("%w: expected uuid got %v", ErrValue, err)
		}

	// TODO split this out into a separate validator
	case EnumValueL:
		u = uuid.UUID(value.(EnumValueL).ID())
	case uuid.UUID:
		u = value.(uuid.UUID)
	case ID:
		u = uuid.UUID(value.(ID))
	default:
		return nil, fmt.Errorf("%w: expected uuid got %T", ErrValue, value)
	}
	return u, nil
}

var uuidValidator = MakeNativeFunction(
	MakeID("60dfeee2-105f-428d-8c10-c4cc3557a40a"),
	"uuid",
	1,
	UUIDFromJSON)

func FloatFromJSON(args []interface{}) (interface{}, error) {
	value := args[0]
	switch value.(type) {
	case float64:
		return value, nil
	case int64:
		return float64(value.(int64)), nil
	case int:
		return float64(value.(int)), nil
	case string:
		out, err := strconv.ParseFloat(value.(string), 64)
		if err == nil {
			return out, nil
		}
	}
	return nil, fmt.Errorf("%w: expected float got %T", ErrValue, value)
}

func TypeFromJSON(args []interface{}) (interface{}, error) {
	value := args[0]
	s, ok := value.(string)
	if !ok {
		return nil, fmt.Errorf("%w: expected string got %T", ErrValue, value)
	}
	return s, ErrNotStored
}

var typeValidator = MakeNativeFunction(
	MakeID("9404dca1-43f5-462f-9916-ad2a22bcb2a7"),
	"type",
	1,
	TypeFromJSON)

var floatValidator = MakeNativeFunction(
	MakeID("83a5f999-00b0-4bc1-879a-434869cf7301"),
	"float",
	1,
	FloatFromJSON)

var Bool = MakeCoreDatatype(
	MakeID("ca05e233-b8a2-4c83-a5c8-87b461c87184"),
	"bool",
	BoolStorage,
	boolValidator,
)

var Int = MakeCoreDatatype(
	MakeID("17cfaaec-7a75-4035-8554-83d8d9194e97"),
	"int",
	IntStorage,
	intValidator,
)

var String = MakeCoreDatatype(
	MakeID("cbab8b98-7ec3-4237-b3e1-eb8bf1112c12"),
	"string",
	StringStorage,
	stringValidator,
)

var UUID = MakeCoreDatatype(
	MakeID("9853fd78-55e6-4dd9-acb9-e04d835eaa42"),
	"uuid",
	UUIDStorage,
	uuidValidator,
)

var Float = MakeCoreDatatype(
	MakeID("72e095f3-d285-47e6-8554-75691c0145e3"),
	"float",
	FloatStorage,
	floatValidator,
)

var Type = MakeCoreDatatype(
	MakeID("e856579f-cebe-44bc-82d7-5a08fde8fe67"),
	"type",
	NotStored,
	typeValidator,
)
