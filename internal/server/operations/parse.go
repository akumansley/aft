package operations

import (
	"awans.org/aft/internal/model"
	"awans.org/aft/internal/server/db"
)

type Parser struct {
	db db.DB
}

func parseAttribute(key string, a model.Attribute, data map[string]interface{}, st interface{}) {
	value := data[key]
	a.SetField(key, value, st)
}

func (p Parser) parseNestedCreate(r model.Relationship, data map[string]interface{}) NestedOperation {
	m := p.db.GetModel(r.Target)
	st := buildStructFromData(m, data)
	nested := []NestedOperation{}
	for k, r := range m.Relationships {
		additionalNested := p.parseRelationship(k, r, data, st)
		nested = append(nested, additionalNested...)
	}
	nestedCreate := NestedCreateOperation{Relationship: r, Struct: st, Nested: nested}
	return nestedCreate
}

func parseNestedConnect(r model.Relationship, data map[string]interface{}) NestedOperation {
	if len(data) != 1 {
		panic("Too many keys in a unique query")
	}
	// this should be a separate method
	var uq UniqueQuery
	for k, v := range data {
		sv := v.(string)
		uq = UniqueQuery{Key: k, Val: sv}
	}
	return NestedConnectOperation{Relationship: r, UniqueQuery: uq}
}

func listify(val interface{}) []interface{} {
	var opList []interface{}
	switch v := val.(type) {
	case map[string]interface{}:
		opList = []interface{}{v}
	case []interface{}:
		opList = v
	default:
		panic("Invalid input")
	}
	return opList
}

func (p Parser) parseRelationship(key string, r model.Relationship, data map[string]interface{}, st interface{}) []NestedOperation {
	nestedOpMap, ok := data[key].(map[string]interface{})
	// Todo: actually check the type -- we don't know we're done here
	if !ok {
		return []NestedOperation{}
	}
	var nested []NestedOperation
	for k, val := range nestedOpMap {
		opList := listify(val)
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
				nestedCreate := p.parseNestedCreate(r, nestedOp)
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

func (p Parser) ParseCreate(modelName string, data map[string]interface{}) CreateOperation {
	m := p.db.GetModel(modelName)
	st := buildStructFromData(m, data)
	nested := []NestedOperation{}
	for k, r := range m.Relationships {
		additionalNested := p.parseRelationship(k, r, data, st)
		nested = append(nested, additionalNested...)
	}
	op := CreateOperation{Struct: st, Nested: nested}
	return op
}
