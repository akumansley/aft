package model

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/ompluscator/dynamic-struct"
	"strings"
)

var typeMap map[AttrType]interface{} = map[AttrType]interface{}{
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
		AttrType: UUID,
	},
	"type": Attribute{
		AttrType: String,
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
		builder.AddField(fieldName, typeMap[sattr.AttrType], fmt.Sprintf(`json:"%v"`, k))
	}

	// later, maybe we can add validate tags
	for k, attr := range m.Attributes {
		fieldName := JsonKeyToFieldName(k)
		builder.AddField(fieldName, typeMap[attr.AttrType], fmt.Sprintf(`json:"%v"`, k))
	}

	for k, rel := range m.Relationships {
		if rel.HasField() {
			relFieldName := JsonKeyToRelFieldName(k)
			builder.AddField(relFieldName, typeMap[UUID], `json:"-"`)
		}
	}

	b := builder.Build()
	memo[modelName] = b
	return b
}
