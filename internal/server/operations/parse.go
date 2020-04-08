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

func parseNestedCreate(r model.Relationship, data map[string]interface{}) NestedOperation {
	m := db.GetModel(r.Target)
	st := buildStructFromData(m, data)
	var nested []NestedOperation
	for k, r := range m.Relationships {
		additionalNested := parseRelationship(k, r, data, st)
		nested = append(nested, additionalNested...)
	}
	nestedCreate := NestedCreateOperation{Relationship: r, Struct: st, Nested: nested}
	return nestedCreate
}

func parseNestedConnect(r model.Relationship, data map[string]interface{}) NestedOperation {
	return NestedConnectOperation{}
}

func parseRelationship(key string, r model.Relationship, data map[string]interface{}, st interface{}) []NestedOperation {
	nestedOpMap, ok := data[key].(map[string]interface{})
	// Todo: actually check the type -- we don't know we're done here
	if !ok {
		return nil
	}
	var nested []NestedOperation
	for k, v := range nestedOpMap {
		// slightly awkward to handle both lists and objects..
		// probably can be refactored
		var opList []interface{}
		obj, isMap := v.(map[string]interface{})
		ls, isLs := v.([]interface{})
		if isMap {
			opList = []interface{}{obj}
		}
		if isLs {
			opList = ls
		}
		for _, op := range opList {
			nestedOp, ok := op.(map[string]interface{})
			if !ok {
				panic("Got list inside of a nested op")
			}
			switch k {
			case "connect":
				nestedConnect := parseNestedConnect(r, nestedOp)
				nested = append(nested, nestedConnect)
			case "create":
				nestedCreate := parseNestedCreate(r, nestedOp)
				nested = append(nested, nestedCreate)
			}
		}
	}

	return nested
}

func buildStructFromData(m model.Model, data map[string]interface{}) interface{} {
	st := model.StructForModel(m).New()
	for k, attr := range m.Attributes {
		parseAttribute(k, attr, data, st)
	}
	return st
}

func ParseCreate(modelName string, data map[string]interface{}) CreateOperation {
	m := db.GetModel(modelName)
	st := buildStructFromData(m, data)
	var nested []NestedOperation
	for k, r := range m.Relationships {
		additionalNested := parseRelationship(k, r, data, st)
		nested = append(nested, additionalNested...)
	}
	op := CreateOperation{Struct: st, Nested: nested}
	return op
}
