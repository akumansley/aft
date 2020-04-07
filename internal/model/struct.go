package model

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/ompluscator/dynamic-struct"
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
	if val, ok := memo[m.Name]; ok {
		return val
	}

	builder := dynamicstruct.NewStruct()

	// later, maybe we can add validate tags
	for k, attr := range m.Attributes {
		fieldName := strings.Title(strings.ToLower(k))
		builder.AddField(fieldName, typeMap[attr.Type], fmt.Sprintf("mapstructure:\"%v\"", k))
	}

	// let's just skip rels for now
	return builder.Build()
}
