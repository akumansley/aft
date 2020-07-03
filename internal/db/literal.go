package db

import (
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
	MarshalDB() ([]Record, []Link)
	GetID() ID
	AsDatatype() Datatype
}

type Link struct {
	from, to ID
	rel      ConcreteRelationshipL
}

func MarshalRecord(v interface{}, lit ModelL) (rec Record) {
	m := lit.AsModel()
	rec = RecordForModel(m)

	attrs, _ := m.Attributes()
	attrMap := map[string]Attribute{}
	for _, a := range attrs {
		attrMap[a.Name()] = a
	}

	vType := reflect.TypeOf(v)
	vVal := reflect.ValueOf(v).Elem()
	for i := 0; i < vType.NumField(); i++ {
		field := vType.Field(i)
		recFieldName := field.Tag.Get("record")
		attr, ok := attrMap[recFieldName]
		if !ok {
			panic("failed to marshal struct to record")
		}
		vIf := vVal.Field(i).Interface()
		err := attr.Set(vIf, rec)
		if err != nil {
			panic("error setting")
		}
	}
	return rec
}
