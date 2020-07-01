package db

import (
	"fmt"
	"github.com/google/uuid"
)

var storageMap map[EnumValue]interface{} = map[EnumValue]interface{}{
	BoolStorage:   false,
	IntStorage:    int64(0),
	StringStorage: "",
	FloatStorage:  0.0,
	UUIDStorage:   uuid.UUID{},
}

func typeCheck(d Datatype, out interface{}) (interface{}, error) {
	switch d.Storage() {
	case BoolStorage:
		return BoolFromJSON(out)
	case IntStorage:
		return IntFromJSON(out)
	case StringStorage:
		return StringFromJSON(out)
	case FloatStorage:
		return FloatFromJSON(out)
	case UUIDStorage:
		return UUIDFromJSON(out)
	}
	return nil, fmt.Errorf("Unrecognized storage for datatype")
}

//Enums
type Enum struct {
	ID   ID
	Name string
}

func (d Enum) FromJSON(arg interface{}) (interface{}, error) {
	//TODO if it's an enum, verify that it's a UUID.
	// I see two options. One we execute some code here that intelligently
	// checks if this enum id matches the enum values from the table.
	// Second option is that I plum through tx into FromJSON. Second seems worse.
	return UUIDFromJSON(arg)
}
