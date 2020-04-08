package operations

import (
	"awans.org/aft/internal/model"
	"awans.org/aft/internal/server/db"
	"fmt"
)

func parseAttribute(key string, a model.Attribute, data map[string]interface{}, st interface{}) {
	value := data[key]
	a.SetField(key, value, st)
}
func parseRelationship(key string, r model.Relationship, data map[string]interface{}, st interface{}) {
}

func ParseCreate(modelName string, data map[string]interface{}) CreateOperation {
	m := db.GetModel(modelName)
	st := model.StructForModel(m).New()

	keysUnused := len(data)
	for k, attr := range m.Attributes {
		parseAttribute(k, attr, data, st)
		keysUnused--
	}
	for k, r := range m.Relationships {
		// may yield nested operations
		parseRelationship(k, r, data, st)
		keysUnused--
	}
	if keysUnused > 0 {
		fmt.Printf("Keysunused: %v\n", keysUnused)
	}
	op := CreateOperation{Struct: st}
	return op
}
