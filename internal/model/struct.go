package model

import (
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

var SystemAttrs = map[string]Attribute{
	"id": Attribute{
		Type: UUID,
	},
	"type": Attribute{
		Type: String,
	},
}

func StructForModel(m Model) dynamicstruct.DynamicStruct {
	modelName := strings.ToLower(m.Name)
	if val, ok := memo[modelName]; ok {
		return val
	}

	builder := dynamicstruct.NewStruct()

	for k, sattr := range SystemAttrs {
		fieldName := JsonKeyToFieldName(k)
		builder.AddField(fieldName, typeMap[sattr.Type], "")
	}

	// later, maybe we can add validate tags
	for k, attr := range m.Attributes {
		fieldName := JsonKeyToFieldName(k)
		builder.AddField(fieldName, typeMap[attr.Type], "")
	}

	for k, rel := range m.Relationships {
		if rel.HasField() {
			relFieldName := JsonKeyToRelFieldName(k)
			builder.AddField(relFieldName, typeMap[UUID], "")
		}
	}

	b := builder.Build()
	memo[modelName] = b
	return b
}
