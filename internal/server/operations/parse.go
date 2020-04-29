package operations

import (
	"awans.org/aft/internal/model"
	"awans.org/aft/internal/server/db"
	"errors"
	"fmt"
	"github.com/google/uuid"
)

var (
	ErrUnusedKeys       = errors.New("unused-keys")
	ErrInvalidStructure = errors.New("invalid-structure")
)

type void struct{}
type set map[string]void

func (s set) String() string {
	var ss []string
	for k := range s {
		ss = append(ss, k)
	}
	return fmt.Sprintf("%v", ss)
}

type Parser struct {
	db db.DB
}

// parseAttribute tries to consume an attribute key from a json map; returns whether the attribute was consumed
func parseAttribute(key string, a model.Attribute, data map[string]interface{}, st interface{}) bool {
	value, ok := data[key]
	if ok {
		a.SetField(key, value, st)
	}
	return ok
}

func (p Parser) parseNestedCreate(r model.Relationship, data map[string]interface{}) (db.NestedOperation, error) {
	unusedKeys := make(set)
	for k := range data {
		unusedKeys[k] = void{}
	}

	m := p.db.GetModel(r.TargetModel)
	st, unusedKeys := buildStructFromData(m, unusedKeys, data)
	nested := []db.NestedOperation{}
	for k, r := range m.Relationships {
		additionalNested, consumed, err := p.parseRelationship(k, r, data, st)
		if err != nil {
			return db.NestedCreateOperation{}, err
		}
		if consumed {
			delete(unusedKeys, k)
		}
		nested = append(nested, additionalNested...)
	}
	if len(unusedKeys) != 0 {
		return db.NestedCreateOperation{}, fmt.Errorf("%w: %v", ErrUnusedKeys, unusedKeys)
	}
	nestedCreate := db.NestedCreateOperation{Relationship: r, Struct: st, Nested: nested}
	return nestedCreate, nil
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

func (p Parser) parseRelationship(key string, r model.Relationship, data map[string]interface{}, st interface{}) ([]db.NestedOperation, bool, error) {
	nestedOpMap, ok := data[key].(map[string]interface{})
	if !ok {
		_, isValue := data[key]
		if !isValue {
			return []db.NestedOperation{}, false, nil
		}

		return []db.NestedOperation{}, false, fmt.Errorf("%w: expected an object, got: %v", ErrInvalidStructure, data)
	}
	var nested []db.NestedOperation
	for k, val := range nestedOpMap {
		opList := listify(val)
		for _, op := range opList {
			nestedOp, ok := op.(map[string]interface{})
			if !ok {
				return nil, false, fmt.Errorf("%w: expected an object, got: %v", ErrInvalidStructure, nestedOp)
			}
			switch k {
			case "connect":
				nestedConnect := parseNestedConnect(r, nestedOp)
				nested = append(nested, nestedConnect)
			case "create":
				nestedCreate, err := p.parseNestedCreate(r, nestedOp)
				if err != nil {
					return nil, false, err
				}
				nested = append(nested, nestedCreate)
			}
		}
	}

	return nested, true, nil
}

func buildStructFromData(m model.Model, keys set, data map[string]interface{}) (interface{}, set) {
	st := model.StructForModel(m).New()
	model.SystemAttrs["type"].SetField("type", m.Name, st)
	for k, attr := range m.Attributes {
		if parseAttribute(k, attr, data, st) {
			delete(keys, k)
		}
	}
	return st, keys
}

func (p Parser) ParseCreate(modelName string, data map[string]interface{}) (op db.CreateOperation, err error) {
	unusedKeys := make(set)
	for k := range data {
		unusedKeys[k] = void{}
	}

	m := p.db.GetModel(modelName)
	st, unusedKeys := buildStructFromData(m, unusedKeys, data)
	nested := []db.NestedOperation{}
	for k, r := range m.Relationships {
		additionalNested, consumed, err := p.parseRelationship(k, r, data, st)
		if err != nil {
			return op, err
		}
		if consumed {
			delete(unusedKeys, k)
		}
		nested = append(nested, additionalNested...)
	}
	if len(unusedKeys) != 0 {
		return op, fmt.Errorf("%w: %v", ErrUnusedKeys, unusedKeys)
	}
	op = db.CreateOperation{Struct: st, Nested: nested}
	return op, err
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
		fieldName = model.JsonKeyToFieldName(k)
		sv, ok := v.(string)
		if !ok {
			panic("non string findOne")
		}
		if m.GetAttributeByJsonName(k).AttrType == model.UUID {
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
	fc := parseFieldCriteria(m, data)
	q.FieldCriteria = fc
	rc := p.parseSingleRelationshipCriteria(m, data)
	q.RelationshipCriteria = rc
	arc := p.parseAggregateRelationshipCriteria(m, data)
	q.AggregateRelationshipCriteria = arc

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

func (p Parser) parseSingleRelationshipCriteria(m model.Model, data map[string]interface{}) []db.RelationshipCriterion {
	var relationshipCriteria []db.RelationshipCriterion
	for k, rel := range m.Relationships {
		if rel.RelType == model.HasOne || rel.RelType == model.BelongsTo {
			if value, ok := data[k]; ok {
				// might need more arguments
				rc := p.parseRelationshipCriterion(rel, value)
				relationshipCriteria = append(relationshipCriteria, rc)
			}
		}
	}
	return relationshipCriteria
}

func (p Parser) parseAggregateRelationshipCriteria(m model.Model, data map[string]interface{}) []db.AggregateRelationshipCriterion {
	var arcs []db.AggregateRelationshipCriterion
	for k, rel := range m.Relationships {
		if rel.RelType == model.HasMany || rel.RelType == model.HasManyAndBelongsToMany {
			if value, ok := data[k]; ok {
				// might need more arguments
				arc := p.parseAggregateRelationshipCriterion(rel, value)
				arcs = append(arcs, arc)
			}
		}
	}
	return arcs
}

func parseFieldCriteria(m model.Model, data map[string]interface{}) []db.FieldCriterion {
	var fieldCriteria []db.FieldCriterion
	for k, attr := range m.Attributes {
		if value, ok := data[k]; ok {
			fc := parseFieldCriterion(k, attr, value)
			fieldCriteria = append(fieldCriteria, fc)
		}
	}
	return fieldCriteria
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

func (p Parser) parseAggregateRelationshipCriterion(r model.Relationship, value interface{}) db.AggregateRelationshipCriterion {
	var arc db.AggregateRelationshipCriterion
	mapValue := value.(map[string]interface{})
	if len(mapValue) > 1 {
		panic("too much data in parseAggregateRel")
	} else if len(mapValue) == 0 {
		panic("empty data in parseAggregateRel")
	}
	var ag db.Aggregation
	for k, v := range mapValue {
		switch k {
		case "some":
			ag = db.Some
		case "none":
			ag = db.None
		case "every":
			ag = db.Every
		default:
			panic("Bad aggregation")
		}
		rc := p.parseRelationshipCriterion(r, v)
		arc = db.AggregateRelationshipCriterion{
			Aggregation:           ag,
			RelationshipCriterion: rc,
		}

	}
	return arc
}

func (p Parser) parseRelationshipCriterion(r model.Relationship, value interface{}) db.RelationshipCriterion {
	mapValue := value.(map[string]interface{})
	m := p.db.GetModel(r.TargetModel)
	fc := parseFieldCriteria(m, mapValue)
	rrc := p.parseSingleRelationshipCriteria(m, mapValue)
	arrc := p.parseAggregateRelationshipCriteria(m, mapValue)
	rc := db.RelationshipCriterion{
		Relationship:                         r,
		RelatedFieldCriteria:                 fc,
		RelatedRelationshipCriteria:          rrc,
		RelatedAggregateRelationshipCriteria: arrc,
	}
	return rc
}

func (p Parser) parseInclusion(r model.Relationship, value interface{}) db.Inclusion {
	if v, ok := value.(bool); ok {
		if v {
			return db.Inclusion{Relationship: r, Query: db.Query{}}
		} else {
			panic("Include specified as false?")
		}
	}
	panic("Include with findMany args not yet implemented")
}

func (p Parser) ParseInclude(modelName string, data map[string]interface{}) db.Include {
	m := p.db.GetModel(modelName)
	var includes []db.Inclusion
	for k, val := range data {
		rel := m.Relationships[k]
		inc := p.parseInclusion(rel, val)
		includes = append(includes, inc)
	}
	io := db.Include{Includes: includes}
	return io
}
