package db

import (
	"fmt"
	"reflect"
)

type Literal interface {
	MarshalDB() ([]Record, []Link)
	ID() ID
}

type AttributeL interface {
	Literal
	Attribute
}

type InterfaceL interface {
	Literal
	Interface
}

type DatatypeL interface {
	Literal
	Datatype
}

type RelationshipL interface {
	Literal
	Relationship
}

type FunctionL interface {
	Literal
	Function
}

type Link struct {
	From, To ID
	Rel      ConcreteRelationshipL
}

func MarshalRecord(v interface{}, m ModelL) (rec Record) {
	rec = RecordForModel(m)

	vType := reflect.TypeOf(v)
	vVal := reflect.ValueOf(v)
	for i := 0; i < vType.NumField(); i++ {
		field := vType.Field(i)
		recFieldName, ok := field.Tag.Lookup("record")
		if !ok {
			continue
		}

		attr, err := m.AttributeByName(recFieldName)
		if err != nil {
			errS := fmt.Errorf("failed to marshal struct to record: %v - %w", recFieldName, err)
			panic(errS)
		}
		vIf := vVal.Field(i).Interface()
		err = attr.Set(rec, vIf)
		if err != nil {
			panic(err)
		}
	}
	return rec
}
