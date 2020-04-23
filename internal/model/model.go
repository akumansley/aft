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

type Attribute struct {
	Type FieldType
}

func (a Attribute) ParseFromJson(value interface{}) interface{} {
	switch a.Type {
	case Int, Enum:
		f, ok := value.(float64)
		i := int64(f)
		if !ok {
			fmt.Printf("Tried setting Int with %v attr %v\n", value, a)
			panic("bad ParseFromJson")
		}
		return i
	case String, Text:
		s, ok := value.(string)
		if !ok {
			fmt.Printf("Tried setting String/Text with %v attr %v\n", value, a)
			panic("bad SetField")
		}
		return s
	case Float:
		f, ok := value.(float64)
		if !ok {
			fmt.Printf("Tried setting float with %v attr %v\n", value, a)
			panic("bad SetField")
		}
		return f
	case UUID:
		uuidString, ok := value.(string)
		if !ok {
			fmt.Printf("Tried setting uuid with %v attr %v\n", value, a)
			panic("bad SetField")
		}
		uuid, err := uuid.Parse(uuidString)
		if err != nil {
			fmt.Printf("couldn't parse uuid")
			panic("bad SetField")
		}
		return uuid
	}
	return nil
}

// arguably this belongs outside of the struct
func (a Attribute) SetField(name string, value interface{}, st interface{}) {
	fieldName := JsonKeyToFieldName(name)
	field := reflect.ValueOf(st).Elem().FieldByName(fieldName)
	parsedValue := a.ParseFromJson(value)
	switch parsedValue.(type) {
	case int64:
		i := parsedValue.(int64)
		field.SetInt(i)
	case string:
		s := parsedValue.(string)
		field.SetString(s)
	case float64:
		f := parsedValue.(float64)
		field.SetFloat(f)
	case uuid.UUID:
		u := parsedValue.(uuid.UUID)
		v := reflect.ValueOf(u)
		field.Set(v)
	}
}

func JsonKeyToFieldName(key string) string {
	return strings.Title(strings.ToLower(key))
}

type RelType int

const (
	HasOne RelType = iota
	BelongsTo
	HasMany
	HasManyAndBelongsToMany
)

type Relationship struct {
	Type        RelType
	TargetModel string
	TargetRel   string
}

func (r Relationship) HasField() bool {
	return r.Type == BelongsTo || r.Type == HasManyAndBelongsToMany
}

func JsonKeyToRelFieldName(key string) string {
	return fmt.Sprintf("%vId", strings.Title(strings.ToLower(key)))
}

type Model struct {
	Type          string
	Id            uuid.UUID
	Name          string
	Attributes    map[string]Attribute
	Relationships map[string]Relationship
}

func (m Model) GetAttributeByJsonName(name string) Attribute {
	a, ok := m.Attributes[name]
	if !ok {
		a, ok = SystemAttrs[name]
	}
	return a
}
