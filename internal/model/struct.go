package model

import (
	"github.com/google/uuid"
	"github.com/ompluscator/dynamic-struct"
	"reflect"
	"strings"
)

var typeMap map[FieldType]interface{} = map[FieldType]interface{}{
	Int:    0,
	String: "",
	Text:   "",
	Float:  0.0,
	Enum:   "",
	UUID:   uuid.UUID{},
}

var memo = map[string]dynamicstruct.DynamicStruct{}

func StructForModel(m Model) dynamicstruct.DynamicStruct {
	modelName := strings.ToLower(m.Name)
	if val, ok := memo[modelName]; ok {
		return val
	}

	builder := dynamicstruct.NewStruct()

	// later, maybe we can add validate tags
	for k, attr := range m.Attributes {
		fieldName := JsonKeyToFieldName(k)
		builder.AddField(fieldName, typeMap[attr.Type], "")
	}

	// let's just skip rels for now
	b := builder.Build()
	memo[modelName] = b
	return b
}

func uUIDDecodeHook(from reflect.Type, to reflect.Type, data interface{}) (interface{}, error) {
	if from == reflect.TypeOf("") && to == reflect.TypeOf(uuid.UUID{}) {
		idString := data.(string)
		u, err := uuid.Parse(idString)
		if err != nil {
			return nil, err
		}
		return u, nil
	}
	return data, nil
}
