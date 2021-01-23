package db

import (
	"fmt"
	"reflect"
	"sort"

	"github.com/google/uuid"
)

type Literal interface {
	MarshalDB(*Builder) ([]Record, []Link)
	ID() ID
	InterfaceID() ID
}

type AttributeL interface {
	Literal
	Load(Tx) Attribute
}

type InterfaceL interface {
	Literal
	Load(Tx) Interface
}

type DatatypeL interface {
	Literal
	Load(Tx) Datatype
}

type RelationshipL interface {
	Literal
	Load(Tx) Relationship
}

type FunctionL interface {
	Literal
	Load(Tx) Function
}

type Link struct {
	From, To ID
	Rel      ConcreteRelationshipL
}

func specFromTaggedLiteral(lit Literal) (s *Spec) {
	s = &Spec{}
	s.InterfaceID = lit.InterfaceID()

	litType := reflect.TypeOf(lit)
	for i := 0; i < litType.NumField(); i++ {
		rField := litType.Field(i)
		fieldName, ok := rField.Tag.Lookup("record")
		if !ok {
			continue
		}
		fieldType := rField.Type
		var storage ID
		switch fieldType {
		case reflect.TypeOf(ID{}):
			storage = UUIDStorage.ID()
		case reflect.TypeOf(""):
			storage = StringStorage.ID()
		case reflect.TypeOf(true):
			storage = BoolStorage.ID()
		case reflect.TypeOf(EnumValueL{}):
			storage = UUIDStorage.ID()
		case reflect.TypeOf(0):
			storage = IntStorage.ID()
		case reflect.TypeOf([]byte{}):
			storage = BytesStorage.ID()
		default:
			err := fmt.Errorf("invalid storage type %v", fieldType)
			panic(err)
		}
		s.Fields = append(s.Fields, Field{Name: fieldName, Storage: storage})
	}
	sort.Slice(s.Fields, func(i, j int) bool {
		return s.Fields[i].Name < s.Fields[j].Name
	})
	return

}

func MarshalRecord(b *Builder, v Literal) (rec Record) {
	rec, err := b.RecordForLiteral(v)
	if err != nil {
		panic(err)
	}

	vType := reflect.TypeOf(v)
	vVal := reflect.ValueOf(v)
	for i := 0; i < vType.NumField(); i++ {
		field := vType.Field(i)
		recFieldName, ok := field.Tag.Lookup("record")
		if !ok {
			continue
		}

		var attrVal interface{}
		attrVal = vVal.Field(i).Interface()
		if ev, ok := attrVal.(EnumValueL); ok {
			attrVal = ev.ID()
		}
		if id, ok := attrVal.(ID); ok {
			attrVal = uuid.UUID(id)
		}
		if i, ok := attrVal.(int); ok {
			attrVal = int64(i)
		}
		rec.Set(recFieldName, attrVal)
	}
	return rec
}
