package operations

import (
	"awans.org/aft/internal/model"
	"awans.org/aft/internal/server/db"
	"errors"
	"fmt"
	"github.com/google/uuid"
)

var (
	ErrParse            = errors.New("parse-error")
	ErrUnusedKeys       = fmt.Errorf("%w: unused keys", ErrParse)
	ErrInvalidModel     = fmt.Errorf("%w: invalid model", ErrParse)
	ErrInvalidStructure = fmt.Errorf("%w: invalid-structure", ErrParse)
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

func (p Parser) parseNestedCreate(r model.Relationship, data map[string]interface{}) (op db.NestedOperation, err error) {
	unusedKeys := make(set)
	for k := range data {
		unusedKeys[k] = void{}
	}

	m, err := p.db.GetModel(r.TargetModel)
	if err != nil {
		return
	}
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

	m, err := p.db.GetModel(modelName)
	if err != nil {
		return op, fmt.Errorf("%w: %v", ErrInvalidModel, modelName)
	}
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

func (p Parser) ParseFindOne(modelName string, data map[string]interface{}) (op db.FindOneOperation, err error) {
	m, err := p.db.GetModel(modelName)
	if err != nil {
		return
	}
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

	op = db.FindOneOperation{
		UniqueQuery: db.UniqueQuery{
			Key: fieldName,
			Val: value,
		},
		ModelName: modelName,
	}
	return op, nil
}

func (p Parser) ParseFindMany(modelName string, data map[string]interface{}) (op db.FindManyOperation, err error) {
	q, err := p.ParseQuery(modelName, data)
	if err != nil {
		return
	}

	op = db.FindManyOperation{
		Query:     q,
		ModelName: modelName,
	}
	return op, nil
}

func (p Parser) parseCompositeQueryList(modelName string, opVal interface{}) (ql []db.Query, err error) {
	opList := opVal.([]interface{})
	for _, opData := range opList {
		opMap := opData.(map[string]interface{})
		var opQ db.Query
		opQ, err = p.ParseQuery(modelName, opMap)
		if err != nil {
			return
		}
		ql = append(ql, opQ)
	}
	return
}

func (p Parser) ParseQuery(modelName string, data map[string]interface{}) (q db.Query, err error) {
	m, err := p.db.GetModel(modelName)
	if err != nil {
		return
	}
	q = db.Query{}
	fc := parseFieldCriteria(m, data)
	q.FieldCriteria = fc
	rc, err := p.parseSingleRelationshipCriteria(m, data)
	if err != nil {
		return
	}
	q.RelationshipCriteria = rc
	arc, err := p.parseAggregateRelationshipCriteria(m, data)
	if err != nil {
		return
	}
	q.AggregateRelationshipCriteria = arc

	if orVal, ok := data["OR"]; ok {
		var orQL []db.Query
		orQL, err = p.parseCompositeQueryList(modelName, orVal)
		if err != nil {
			return
		}
		q.Or = orQL
	}
	if andVal, ok := data["AND"]; ok {
		var andQL []db.Query
		andQL, err = p.parseCompositeQueryList(modelName, andVal)
		if err != nil {
			return
		}
		q.And = andQL
	}
	if notVal, ok := data["NOT"]; ok {
		var notQL []db.Query
		notQL, err = p.parseCompositeQueryList(modelName, notVal)
		if err != nil {
			return
		}
		q.Not = notQL
	}
	return
}

func (p Parser) parseSingleRelationshipCriteria(m model.Model, data map[string]interface{}) (rcl []db.RelationshipCriterion, err error) {
	for k, rel := range m.Relationships {
		if rel.RelType == model.HasOne || rel.RelType == model.BelongsTo {
			if value, ok := data[k]; ok {
				var rc db.RelationshipCriterion
				rc, err = p.parseRelationshipCriterion(rel, value)
				if err != nil {
					return
				}
				rcl = append(rcl, rc)
			}
		}
	}
	return rcl, nil
}

func (p Parser) parseAggregateRelationshipCriteria(m model.Model, data map[string]interface{}) (arcl []db.AggregateRelationshipCriterion, err error) {
	for k, rel := range m.Relationships {
		if rel.RelType == model.HasMany || rel.RelType == model.HasManyAndBelongsToMany {
			if value, ok := data[k]; ok {
				var arc db.AggregateRelationshipCriterion
				arc, err = p.parseAggregateRelationshipCriterion(rel, value)
				if err != nil {
					return
				}
				arcl = append(arcl, arc)
			}
		}
	}
	return arcl, nil
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

func (p Parser) parseAggregateRelationshipCriterion(r model.Relationship, value interface{}) (arc db.AggregateRelationshipCriterion, err error) {
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
		var rc db.RelationshipCriterion
		rc, err = p.parseRelationshipCriterion(r, v)
		if err != nil {
			return
		}
		arc = db.AggregateRelationshipCriterion{
			Aggregation:           ag,
			RelationshipCriterion: rc,
		}
	}
	return
}

func (p Parser) parseRelationshipCriterion(r model.Relationship, value interface{}) (rc db.RelationshipCriterion, err error) {
	mapValue := value.(map[string]interface{})
	m, err := p.db.GetModel(r.TargetModel)
	if err != nil {
		return
	}
	fc := parseFieldCriteria(m, mapValue)
	rrc, err := p.parseSingleRelationshipCriteria(m, mapValue)
	if err != nil {
		return
	}
	arrc, err := p.parseAggregateRelationshipCriteria(m, mapValue)
	if err != nil {
		return
	}
	rc = db.RelationshipCriterion{
		Relationship:                         r,
		RelatedFieldCriteria:                 fc,
		RelatedRelationshipCriteria:          rrc,
		RelatedAggregateRelationshipCriteria: arrc,
	}
	return
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

func (p Parser) ParseInclude(modelName string, data map[string]interface{}) (i db.Include, err error) {
	m, err := p.db.GetModel(modelName)
	if err != nil {
		return
	}
	var includes []db.Inclusion
	for k, val := range data {
		rel := m.Relationships[k]
		inc := p.parseInclusion(rel, val)
		includes = append(includes, inc)
	}
	i = db.Include{Includes: includes}
	return
}
