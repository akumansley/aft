package db

import (
	"fmt"
	"github.com/google/uuid"
	"strconv"
)

var (
	ErrValue = fmt.Errorf("invalid value for type")
)

func BoolFromJSON(value interface{}) (interface{}, error) {
	b, ok := value.(bool)
	if !ok {
		return nil, fmt.Errorf("%w: expected bool got %T", ErrValue, value)
	}
	return b, nil
}

var boolValidator = MakeNativeFunction(
	MakeID("8e806967-c462-47af-8756-48674537a909"),
	"bool",
	FromJSON,
	BoolFromJSON)

func IntFromJSON(value interface{}) (interface{}, error) {
	switch value.(type) {
	case int64:
		return value, nil
	case int:
		return int64(value.(int)), nil
	case float64:
		return int64(value.(float64)), nil
	case string:
		out, err := strconv.Atoi(value.(string))
		if err == nil {
			return out, nil
		}
	}
	return nil, fmt.Errorf("%w: expected int got %T", ErrValue, value)

}

var intValidator = MakeNativeFunction(
	MakeID("a1cf1c16-040d-482c-92ae-92d59dbad46c"),
	"int",
	FromJSON,
	IntFromJSON)

func StringFromJSON(value interface{}) (interface{}, error) {
	s, ok := value.(string)
	if !ok {
		return nil, fmt.Errorf("%w: expected string got %T", ErrValue, value)
	}
	return s, nil
}

var stringValidator = MakeNativeFunction(
	MakeID("aaeccd14-e69f-4561-91ef-5a8a75b0b498"),
	"string",
	FromJSON,
	StringFromJSON)

func UUIDFromJSON(value interface{}) (interface{}, error) {
	var u uuid.UUID
	var err error
	switch value.(type) {
	case string:
		u, err = uuid.Parse(value.(string))
		if err != nil {
			return nil, fmt.Errorf("%w: expected uuid got %v", ErrValue, err)
		}
	case uuid.UUID:
		u = value.(uuid.UUID)
	default:
		return nil, fmt.Errorf("%w: expected uuid got %T", ErrValue, value)
	}
	return u, nil
}

var uuidValidator = MakeNativeFunction(
	MakeID("60dfeee2-105f-428d-8c10-c4cc3557a40a"),
	"uuid",
	FromJSON,
	UUIDFromJSON)

func FloatFromJSON(value interface{}) (interface{}, error) {
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

var floatValidator = MakeNativeFunction(
	MakeID("83a5f999-00b0-4bc1-879a-434869cf7301"),
	"float",
	FromJSON,
	FloatFromJSON)

var Bool = CoreDatatypeL{
	ID:        MakeID("ca05e233-b8a2-4c83-a5c8-87b461c87184"),
	Name:      "bool",
	Validator: boolValidator,
	StoredAs:  BoolStorage,
}

var Int = CoreDatatypeL{
	ID:        MakeID("17cfaaec-7a75-4035-8554-83d8d9194e97"),
	Name:      "int",
	Validator: intValidator,
	StoredAs:  IntStorage,
}

var String = CoreDatatypeL{
	ID:        MakeID("cbab8b98-7ec3-4237-b3e1-eb8bf1112c12"),
	Name:      "string",
	Validator: stringValidator,
	StoredAs:  StringStorage,
}

var UUID = CoreDatatypeL{
	ID:        MakeID("9853fd78-55e6-4dd9-acb9-e04d835eaa42"),
	Name:      "uuid",
	Validator: uuidValidator,
	StoredAs:  UUIDStorage,
}

var Float = CoreDatatypeL{
	ID:        MakeID("72e095f3-d285-47e6-8554-75691c0145e3"),
	Name:      "float",
	Validator: floatValidator,
	StoredAs:  FloatStorage,
}
