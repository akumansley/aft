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

func (p Parser) parseNestedCreate(r model.Relationship, data map[string]interface{}) db.NestedOperation {
	m := p.db.GetModel(r.TargetModel)
	st := buildStructFromData(m, data)
	nested := []db.NestedOperation{}
	for k, r := range m.Relationships {
		additionalNested := p.parseRelationship(k, r, data, st)
		nested = append(nested, additionalNested...)
	}
	nestedCreate := db.NestedCreateOperation{Relationship: r, Struct: st, Nested: nested}
	return nestedCreate
}

func parseNestedConnect(r model.Relationship, data map[string]interface{}) db.NestedOperation {
	if len(data) != 1 {
		panic("Too many keys in a unique query")
	}
	// this should be a separate method
	var uq db.UniqueQuery
	for k, v := range data {
		sv := v.(string)
		uq = db.UniqueQuery{Key: k, Val: sv}
	}
	return db.NestedConnectOperation{Relationship: r, UniqueQuery: uq}
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

func (p Parser) parseRelationship(key string, r model.Relationship, data map[string]interface{}, st interface{}) []db.NestedOperation {
	nestedOpMap, ok := data[key].(map[string]interface{})
	// Todo: actually check the type -- we don't know we're done here
	if !ok {
		return []db.NestedOperation{}
	}
	var nested []db.NestedOperation
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
	for k, sattr := range model.SystemAttrs {
		parseAttribute(k, sattr, data, st)
	}
	for k, attr := range m.Attributes {
		parseAttribute(k, attr, data, st)
	}
	return st
}

func (p Parser) ParseCreate(modelName string, data map[string]interface{}) db.CreateOperation {
	m := p.db.GetModel(modelName)
	st := buildStructFromData(m, data)
	nested := []db.NestedOperation{}
	for k, r := range m.Relationships {
		additionalNested := p.parseRelationship(k, r, data, st)
		nested = append(nested, additionalNested...)
	}
	op := db.CreateOperation{Struct: st, Nested: nested}
	return op
}
