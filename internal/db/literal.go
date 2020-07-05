package db

import (
	"fmt"
	"reflect"
)

type Literal interface {
	MarshalDB() ([]Record, []Link)
	GetID() ID
}

type AttributeL interface {
	MarshalDB() ([]Record, []Link)
	GetID() ID
	AsAttribute() Attribute
}

type DatatypeL interface {
	Literal
	Datatype
}

type Link struct {
	from, to ID
	rel      ConcreteRelationshipL
}

func MarshalRecord(v interface{}, lit ModelL) (rec Record) {
	m := lit.AsModel()
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
		err = attr.Set(vIf, rec)
		if err != nil {
			panic(err)
		}
	}
	return rec
}
