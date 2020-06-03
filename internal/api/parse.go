package api

import (
	"awans.org/aft/internal/db"
	"errors"
	"fmt"
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
	tx db.Tx
}

// parseAttribute tries to consume an attribute key from a json map; returns whether the attribute was consumed
func parseAttribute(key string, data map[string]interface{}, rec db.Record) (ok bool, err error) {
	value, ok := data[key]
	if ok {
		err := rec.Set(key, value)
		if err != nil {
			return false, err
		}
	}
	return ok, nil
}

func (p Parser) parseNestedCreate(parentBinding db.Binding, data map[string]interface{}) (op NestedOperation, err error) {
	unusedKeys := make(set)
	for k := range data {
		unusedKeys[k] = void{}
	}

	targetModel, err := p.tx.GetModelByID(parentBinding.Dual().ModelID())
	if err != nil {
		return
	}
	rec, unusedKeys, err := buildRecordFromData(targetModel, unusedKeys, data)
	if err != nil {
		return
	}
	nested := []NestedOperation{}
	for _, b := range targetModel.Bindings() {
		additionalNested, consumed, err := p.parseRelationship(b, data)
		if err != nil {
			return NestedCreateOperation{}, err
		}
		if consumed {
			delete(unusedKeys, b.Name())
		}
		nested = append(nested, additionalNested...)
	}
	if len(unusedKeys) != 0 {
		return NestedCreateOperation{}, fmt.Errorf("%w: %v", ErrUnusedKeys, unusedKeys)
	}
	nestedCreate := NestedCreateOperation{Binding: parentBinding, Record: rec, Nested: nested}
	return nestedCreate, nil
}

func (p Parser) parseNestedConnect(parentBinding db.Binding, data map[string]interface{}) (op NestedConnectOperation, err error) {
	if len(data) != 1 {
		panic("Too many keys in a unique query")
	}
	m, err := p.tx.GetModelByID(parentBinding.Dual().ModelID())
	if err != nil {
		return
	}
	// this should be a separate method
	var uq UniqueQuery
	for k, v := range data {
		var val interface{}
		d := m.AttributeByName(k).Datatype
		val, err = db.CallFunc(d.FromJSON, v, d.StorageType)
		if err != nil {
			return op, fmt.Errorf("error parsing %v %v: %w", m.Name, k, err)
		}
		uq = UniqueQuery{Key: k, Val: val}
	}
	return NestedConnectOperation{Binding: parentBinding, UniqueQuery: uq}, nil
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

func (p Parser) parseRelationship(b db.Binding, data map[string]interface{}) ([]NestedOperation, bool, error) {
	// refactor this
	nestedOpMap, ok := data[b.Name()].(map[string]interface{})
	if !ok {
		_, isValue := data[b.Name()]
		if !isValue {
			return []NestedOperation{}, false, nil
		}

		return []NestedOperation{}, false, fmt.Errorf("%w: expected an object, got: %v", ErrInvalidStructure, data)
	}
	var nested []NestedOperation
	for k, val := range nestedOpMap {
		opList := listify(val)
		for _, op := range opList {
			nestedOp, ok := op.(map[string]interface{})
			if !ok {
				return nil, false, fmt.Errorf("%w: expected an object, got: %v", ErrInvalidStructure, nestedOp)
			}
			switch k {
			case "connect":
				nestedConnect, err := p.parseNestedConnect(b, nestedOp)
				if err != nil {
					return nil, false, err
				}
				nested = append(nested, nestedConnect)
			case "create":
				nestedCreate, err := p.parseNestedCreate(b, nestedOp)
				if err != nil {
					return nil, false, err
				}
				nested = append(nested, nestedCreate)
			}
		}
	}

	return nested, true, nil
}

func buildRecordFromData(m db.Model, keys set, data map[string]interface{}) (db.Record, set, error) {
	rec := db.RecordForModel(m)
	for k := range m.Attributes {
		ok, err := parseAttribute(k, data, rec)
		if err != nil {
			return nil, nil, err
		}
		if ok {
			delete(keys, k)
		}
	}
	return rec, keys, nil
}

func (p Parser) ParseCreate(modelName string, data map[string]interface{}) (op CreateOperation, err error) {
	unusedKeys := make(set)
	for k := range data {
		unusedKeys[k] = void{}
	}

	m, err := p.tx.GetModel(modelName)
	if err != nil {
		return op, fmt.Errorf("%w: %v", ErrInvalidModel, modelName)
	}
	rec, unusedKeys, err := buildRecordFromData(m, unusedKeys, data)
	if err != nil {
		return op, fmt.Errorf("%w: %v", ErrParse, err)
	}
	nested := []NestedOperation{}
	for _, b := range m.Bindings() {
		additionalNested, consumed, err := p.parseRelationship(b, data)
		if err != nil {
			return op, err
		}
		if consumed {
			delete(unusedKeys, b.Name())
		}
		nested = append(nested, additionalNested...)
	}
	if len(unusedKeys) != 0 {
		return op, fmt.Errorf("%w: %v", ErrUnusedKeys, unusedKeys)
	}
	op = CreateOperation{Record: rec, Nested: nested}
	return op, err
}

func (p Parser) ParseFindOne(modelName string, data map[string]interface{}) (op FindOneOperation, err error) {
	m, err := p.tx.GetModel(modelName)
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
		d := m.AttributeByName(k).Datatype
		fieldName = db.JSONKeyToFieldName(k)
		value, err = db.CallFunc(d.FromJSON, v, d.StorageType)
		if err != nil {
			return
		}
	}

	op = FindOneOperation{
		UniqueQuery: UniqueQuery{
			Key: fieldName,
			Val: value,
		},
		ModelName: modelName,
	}
	return
}

func (p Parser) ParseFindMany(modelName string, data map[string]interface{}) (op FindManyOperation, err error) {
	q, err := p.ParseQuery(modelName, data)
	if err != nil {
		return
	}

	op = FindManyOperation{
		Query:     q,
		ModelName: modelName,
	}
	return op, nil
}

func (p Parser) parseCompositeQueryList(modelName string, opVal interface{}) (ql []Query, err error) {
	opList := opVal.([]interface{})
	for _, opData := range opList {
		opMap := opData.(map[string]interface{})
		var opQ Query
		opQ, err = p.ParseQuery(modelName, opMap)
		if err != nil {
			return
		}
		ql = append(ql, opQ)
	}
	return
}

func (p Parser) ParseQuery(modelName string, data map[string]interface{}) (q Query, err error) {
	m, err := p.tx.GetModel(modelName)
	if err != nil {
		return
	}
	q = Query{}
	fc, err := parseFieldCriteria(m, data)
	if err != nil {
		return
	}
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
		var orQL []Query
		orQL, err = p.parseCompositeQueryList(modelName, orVal)
		if err != nil {
			return
		}
		q.Or = orQL
	}
	if andVal, ok := data["AND"]; ok {
		var andQL []Query
		andQL, err = p.parseCompositeQueryList(modelName, andVal)
		if err != nil {
			return
		}
		q.And = andQL
	}
	if notVal, ok := data["NOT"]; ok {
		var notQL []Query
		notQL, err = p.parseCompositeQueryList(modelName, notVal)
		if err != nil {
			return
		}
		q.Not = notQL
	}
	return
}

func (p Parser) parseSingleRelationshipCriteria(m db.Model, data map[string]interface{}) (rcl []RelationshipCriterion, err error) {
	for _, b := range m.Bindings() {
		if b.RelType() == db.HasOne || b.RelType() == db.BelongsTo {
			if value, ok := data[b.Name()]; ok {
				var rc RelationshipCriterion
				rc, err = p.parseRelationshipCriterion(b, value)
				if err != nil {
					return
				}
				rcl = append(rcl, rc)
			}
		}
	}
	return rcl, nil
}

func (p Parser) parseAggregateRelationshipCriteria(m db.Model, data map[string]interface{}) (arcl []AggregateRelationshipCriterion, err error) {
	for _, b := range m.Bindings() {
		if b.RelType() == db.HasMany || b.RelType() == db.HasManyAndBelongsToMany {
			if value, ok := data[b.Name()]; ok {
				var arc AggregateRelationshipCriterion
				arc, err = p.parseAggregateRelationshipCriterion(b, value)
				if err != nil {
					return
				}
				arcl = append(arcl, arc)
			}
		}
	}
	return arcl, nil
}

func parseFieldCriteria(m db.Model, data map[string]interface{}) (fieldCriteria []FieldCriterion, err error) {
	for k, attr := range m.Attributes {
		if value, ok := data[k]; ok {
			var fc FieldCriterion
			fc, err = parseFieldCriterion(k, attr, value)
			fieldCriteria = append(fieldCriteria, fc)
		}
	}
	return
}

func parseFieldCriterion(key string, a db.Attribute, value interface{}) (fc FieldCriterion, err error) {
	fieldName := db.JSONKeyToFieldName(key)
	parsedValue, err := db.CallFunc(a.Datatype.FromJSON, value, a.Datatype.StorageType)
	fc = FieldCriterion{
		// TODO handle function values like {startsWith}
		Key: fieldName,
		Val: parsedValue,
	}
	return
}

func (p Parser) parseAggregateRelationshipCriterion(b db.Binding, value interface{}) (arc AggregateRelationshipCriterion, err error) {
	mapValue := value.(map[string]interface{})
	if len(mapValue) > 1 {
		panic("too much data in parseAggregateRel")
	} else if len(mapValue) == 0 {
		panic("empty data in parseAggregateRel")
	}
	var ag Aggregation
	for k, v := range mapValue {
		switch k {
		case "some":
			ag = Some
		case "none":
			ag = None
		case "every":
			ag = Every
		default:
			panic("Bad aggregation")
		}
		var rc RelationshipCriterion
		rc, err = p.parseRelationshipCriterion(b, v)
		if err != nil {
			return
		}
		arc = AggregateRelationshipCriterion{
			Aggregation:           ag,
			RelationshipCriterion: rc,
		}
	}
	return
}

func (p Parser) parseRelationshipCriterion(b db.Binding, value interface{}) (rc RelationshipCriterion, err error) {
	mapValue := value.(map[string]interface{})
	m, err := p.tx.GetModelByID(b.Dual().ModelID())
	if err != nil {
		return
	}
	fc, err := parseFieldCriteria(m, mapValue)
	if err != nil {
		return
	}
	rrc, err := p.parseSingleRelationshipCriteria(m, mapValue)
	if err != nil {
		return
	}
	arrc, err := p.parseAggregateRelationshipCriteria(m, mapValue)
	if err != nil {
		return
	}
	rc = RelationshipCriterion{
		Binding:                              b,
		RelatedFieldCriteria:                 fc,
		RelatedRelationshipCriteria:          rrc,
		RelatedAggregateRelationshipCriteria: arrc,
	}
	return
}

func (p Parser) parseInclusion(b db.Binding, value interface{}) Inclusion {
	if v, ok := value.(bool); ok {
		if v {
			return Inclusion{Binding: b, Query: Query{}}
		} else {
			panic("Include specified as false?")
		}
	}
	panic("Include with findMany args not yet implemented")
}

func (p Parser) ParseInclude(modelName string, data map[string]interface{}) (i Include, err error) {
	m, err := p.tx.GetModel(modelName)
	if err != nil {
		return
	}
	var includes []Inclusion
	for k, val := range data {
		var b db.Binding
		b, err = m.GetBinding(k)
		if err != nil {
			return
		}
		inc := p.parseInclusion(b, val)
		includes = append(includes, inc)
	}
	i = Include{Includes: includes}
	return
}
