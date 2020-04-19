package operations

import (
	"awans.org/aft/internal/model"
	"awans.org/aft/internal/server/db"
	"fmt"
	"github.com/google/uuid"
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

func (p Parser) ParseFindOne(modelName string, data map[string]interface{}) db.FindOneOperation {
	m := p.db.GetModel(modelName)
	var fieldName string
	var value interface{}

	if len(data) > 1 {
		panic("too much data in findOne")
	} else if len(data) == 0 {
		panic("empty data in findOne")
	}

	for k, v := range data {
		fmt.Printf("%v:%v\n", k, v)
		fieldName = model.JsonKeyToFieldName(k)
		sv, ok := v.(string)
		if !ok {
			panic("non string findOne")
		}
		if m.GetAttributeByJsonName(k).Type == model.UUID {
			uValue, err := uuid.Parse(sv)
			if err != nil {
				panic("uuid failed to parse")
			}
			value = uValue
		} else {
			value = sv
		}
	}

	op := db.FindOneOperation{
		UniqueQuery: db.UniqueQuery{
			Key: fieldName,
			Val: value,
		},
		ModelName: modelName,
	}
	return op
}

func (p Parser) ParseFindMany(modelName string, data map[string]interface{}) db.FindManyOperation {
	q := p.ParseQuery(modelName, data)

	op := db.FindManyOperation{
		Query:     q,
		ModelName: modelName,
	}
	return op
}

func (p Parser) parseCompositeQueryList(modelName string, opVal interface{}) []db.Query {
	var opQueries []db.Query
	opList := opVal.([]interface{})
	for _, opData := range opList {
		opMap := opData.(map[string]interface{})
		opQ := p.ParseQuery(modelName, opMap)
		opQueries = append(opQueries, opQ)
	}
	return opQueries
}

func (p Parser) ParseQuery(modelName string, data map[string]interface{}) db.Query {
	m := p.db.GetModel(modelName)
	q := db.Query{}
	var fieldCriteria []db.FieldCriterion
	for k, attr := range m.Attributes {
		if value, ok := data[k]; ok {
			fc := parseFieldCriterion(k, attr, value)
			fieldCriteria = append(fieldCriteria, fc)
		}
	}
	q.FieldCriteria = fieldCriteria

	var relationshipCriteria []db.RelationshipCriterion
	for k, rel := range m.Relationships {
		if value, ok := data[k]; ok {
			// might need more arguments
			rc := parseRelationshipCriterion(rel, value)
			relationshipCriteria = append(relationshipCriteria, rc)
		}
	}
	q.RelationshipCriteria = relationshipCriteria

	if orVal, ok := data["OR"]; ok {
		orQL := p.parseCompositeQueryList(modelName, orVal)
		q.Or = orQL
	}
	if andVal, ok := data["AND"]; ok {
		andQL := p.parseCompositeQueryList(modelName, andVal)
		q.And = andQL
	}
	if notVal, ok := data["NOT"]; ok {
		notQL := p.parseCompositeQueryList(modelName, notVal)
		q.Not = notQL
	}
	return q
}

func parseFieldCriterion(key string, a model.Attribute, value interface{}) db.FieldCriterion {
	fieldName := model.JsonKeyToFieldName(key)
	parsedValue := a.ParseFromJson(value)
	fc := db.FieldCriterion{
		// TODO handle function values like {startsWith}
		Key: fieldName,
		Val: parsedValue,
	}
	return fc
}

func parseRelationshipCriterion(r model.Relationship, value interface{}) db.RelationshipCriterion {
	// TODO handle to-one rels
	// should be similar to parsequery
	// TODO handle to-many rels
	// some/none/all

	return db.RelationshipCriterion{}
}
