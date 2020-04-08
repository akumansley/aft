package model

import (
	"fmt"
	"github.com/google/uuid"
	"reflect"
	"strings"
)

type FieldType int

const (
	Int FieldType = iota
	String
	Text
	Float
	Enum
	UUID
)

type RelType int

const (
	HasOne RelType = iota
	BelongsTo
	HasMany
	HasManyAndBelongsToMany
)

type Attribute struct {
	Type FieldType
}

// arguably this belongs outside of the struct
func (a Attribute) SetField(name string, value interface{}, st interface{}) {
	fieldName := JsonKeyToFieldName(name)
	field := reflect.ValueOf(st).Elem().FieldByName(fieldName)
	switch a.Type {
	case Int, Enum:
		f, ok := value.(float64)
		i := int64(f)
		if !ok {
			panic("SetField was bogus")
		}
		field.SetInt(i)
	case String, Text:
		s, ok := value.(string)
		if !ok {
			panic("SetField was bogus")
		}
		field.SetString(s)
	case Float:
		f, ok := value.(float64)
		if !ok {
			panic("setfield was bogus")
		}
		field.SetFloat(f)
	case UUID:
		uuidString, ok := value.(string)
		if !ok {
			fmt.Printf("value is %v", value)
			panic("setfield was bogus ")
		}
		uuid, err := uuid.Parse(uuidString)
		if err != nil {
			panic("couldn't parse uuid")
		}
		v := reflect.ValueOf(uuid)
		field.Set(v)
	}
}

func JsonKeyToFieldName(key string) string {
	return strings.Title(strings.ToLower(key))
}

type Relationship struct {
	Type   RelType
	Target string
}

type Model struct {
	Name          string `boldholdIndex:"Name"`
	Attributes    map[string]Attribute
	Relationships map[string]Relationship
}
