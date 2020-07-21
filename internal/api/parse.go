package api

import (
	"awans.org/aft/internal/db"
	"errors"
	"fmt"
)

var (
	ErrParse               = errors.New("parse-error")
	ErrUnusedKeys          = fmt.Errorf("%w: unused keys", ErrParse)
	ErrInvalidModel        = fmt.Errorf("%w: invalid model", ErrParse)
	ErrInvalidRelationship = fmt.Errorf("%w: invalid relationship", ErrParse)
	ErrInvalidStructure    = fmt.Errorf("%w: invalid-structure", ErrParse)
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
func parseAttribute(attr db.Attribute, data map[string]interface{}, rec db.Record) (ok bool, err error) {
	key := attr.Name()
	value, ok := data[key]
	if ok {
		err := attr.Set(rec, value)
		if err != nil {
			return false, err
		}
	}
	return ok, nil
}

func (p Parser) parseNestedCreate(rel db.Relationship, data map[string]interface{}) (op NestedOperation, err error) {
	unusedKeys := make(set)
	for k := range data {
		unusedKeys[k] = void{}
	}

	targetModel := rel.Target()
	rec, unusedKeys, err := buildRecordFromData(targetModel, unusedKeys, data)
	if err != nil {
		return
	}
	nested := []NestedOperation{}
	rels, err := targetModel.Relationships()
	if err != nil {
		return
	}
	for _, r := range rels {
		additionalNested, consumed, err := p.parseRelationship(r, data)
		if err != nil {
			return NestedCreateOperation{}, err
		}
		if consumed {
			delete(unusedKeys, r.Name())
		}
		nested = append(nested, additionalNested...)
	}
	if len(unusedKeys) != 0 {
		return NestedCreateOperation{}, fmt.Errorf("%w: %v", ErrUnusedKeys, unusedKeys)
	}
	nestedCreate := NestedCreateOperation{Relationship: rel, Record: rec, Nested: nested}
	return nestedCreate, nil
}

func (p Parser) parseNestedConnect(rel db.Relationship, data map[string]interface{}) (op NestedConnectOperation, err error) {
	if len(data) != 1 {
		panic("Too many keys in a unique query")
	}
	m := rel.Target()

	// this should be a separate method
	var uq UniqueQuery
	for k, v := range data {
		var val interface{}
		a, err := m.AttributeByName(k)
		d := a.Datatype()

		if err != nil {
			return op, fmt.Errorf("error parsing %v %v: %w", m.Name(), k, err)
		}
		f, err := d.FromJSON()
		if err != nil {
			return op, fmt.Errorf("error parsing %v %v: %w", m.Name(), k, err)
		}
		val, err = f.Call(v)
		if err != nil {
			return op, fmt.Errorf("error parsing %v %v: %w", m.Name(), k, err)
		}
		uq = UniqueQuery{Key: k, Val: val}
	}
	return NestedConnectOperation{Relationship: rel, UniqueQuery: uq}, nil
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

func (p Parser) parseRelationship(r db.Relationship, data map[string]interface{}) ([]NestedOperation, bool, error) {
	// refactor this
	nestedOpMap, ok := data[r.Name()].(map[string]interface{})
	if !ok {
		_, isValue := data[r.Name()]
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
				nestedConnect, err := p.parseNestedConnect(r, nestedOp)
				if err != nil {
					return nil, false, err
				}
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

func buildRecordFromData(m db.Interface, keys set, data map[string]interface{}) (db.Record, set, error) {
	rec := db.NewRecord(m)
	attrs, err := m.Attributes()
	if err != nil {
		return nil, keys, err
	}
	for _, a := range attrs {
		ok, err := parseAttribute(a, data, rec)
		if err != nil {
			return nil, nil, err
		}
		if ok {
			delete(keys, a.Name())
		}
	}
	return rec, keys, nil
}

func updateRecordFromData(oldRec db.Record, keys set, data map[string]interface{}) (db.Record, set, error) {
	newRec := oldRec.DeepCopy()
	attrs, err := oldRec.Interface().Attributes()
	if err != nil {
		return nil, keys, err
	}
	for _, a := range attrs {
		ok, err := parseAttribute(a, data, newRec)
		if err != nil {
			return nil, nil, err
		}
		if ok {
			delete(keys, a.Name())
		}
	}
	return newRec, keys, nil
}

func (p Parser) ParseCreate(modelName string, data map[string]interface{}) (op CreateOperation, err error) {
	unusedKeys := make(set)
	for k := range data {
		unusedKeys[k] = void{}
	}

	m, err := p.tx.Schema().GetInterface(modelName)
	if err != nil {
		return op, fmt.Errorf("%w: %v", ErrInvalidModel, modelName)
	}
	rec, unusedKeys, err := buildRecordFromData(m, unusedKeys, data)
	if err != nil {
		return op, fmt.Errorf("%w: %v", ErrParse, err)
	}
	nested := []NestedOperation{}
	rels, err := m.Relationships()
	if err != nil {
		return
	}
	for _, r := range rels {
		additionalNested, consumed, err := p.parseRelationship(r, data)
		if err != nil {
			return op, err
		}
		if consumed {
			delete(unusedKeys, r.Name())
		}
		nested = append(nested, additionalNested...)
	}
	if len(unusedKeys) != 0 {
		return op, fmt.Errorf("%w: %v", ErrUnusedKeys, unusedKeys)
	}
	op = CreateOperation{Record: rec, Nested: nested}
	return op, err
}

func (p Parser) ParseUpdate(oldRec db.Record, data map[string]interface{}) (op UpdateOperation, err error) {
	unusedKeys := make(set)
	for k := range data {
		unusedKeys[k] = void{}
	}

	newRec, unusedKeys, err := updateRecordFromData(oldRec, unusedKeys, data)
	if len(unusedKeys) != 0 {
		return op, fmt.Errorf("%w: %v", ErrUnusedKeys, unusedKeys)
	}
	op = UpdateOperation{Old: oldRec, New: newRec}
	return op, err
}

func (p Parser) ParseUpdateMany(oldRecs []db.Record, data map[string]interface{}) (op UpdateManyOperation, err error) {
	unusedKeys := make(set)
	for k := range data {
		unusedKeys[k] = void{}
	}
	var newRecs []db.Record
	for _, oldRec := range oldRecs {
		newRec, unusedKeys, err := updateRecordFromData(oldRec, unusedKeys, data)
		if err != nil {
			return op, err
		}
		newRecs = append(newRecs, newRec)
		if len(unusedKeys) != 0 {
			return op, fmt.Errorf("%w: %v", ErrUnusedKeys, unusedKeys)
		}
	}
	op = UpdateManyOperation{Old: oldRecs, New: newRecs}
	return op, err
}

func (p Parser) ParseFindOne(modelName string, data map[string]interface{}) (op FindOneOperation, err error) {
	m, err := p.tx.Schema().GetInterface(modelName)
	if err != nil {
		return
	}
	var fieldName string
	var value interface{}

	if len(data) > 1 {
		return op, fmt.Errorf("%w: %v", ErrInvalidStructure, data)
	} else if len(data) == 0 {
		return op, fmt.Errorf("%w: %v", ErrInvalidStructure, data)
	}

	for k, v := range data {
		fieldName = db.JSONKeyToFieldName(k)
		var a db.Attribute
		a, err = m.AttributeByName(k)
		if err != nil {
			return
		}
		d := a.Datatype()
		var f db.Function
		f, err = d.FromJSON()
		if err != nil {
			return
		}

		value, err = f.Call(v)
		if err != nil {
			return
		}
	}

	op = FindOneOperation{
		UniqueQuery: UniqueQuery{
			Key: fieldName,
			Val: value,
		},
		ModelID: m.ID(),
	}
	return
}

func (p Parser) ParseFindMany(modelName string, data map[string]interface{}) (op FindManyOperation, err error) {
	m, err := p.tx.Schema().GetInterface(modelName)
	if err != nil {
		return
	}
	q, err := p.ParseWhere(modelName, data)
	if err != nil {
		return
	}

	op = FindManyOperation{
		Where:   q,
		ModelID: m.ID(),
	}
	return op, nil
}

func (p Parser) parseCompositeQueryList(modelName string, opVal interface{}) (ql []Where, err error) {
	opList := opVal.([]interface{})
	for _, opData := range opList {
		opMap := opData.(map[string]interface{})
		var opQ Where
		opQ, err = p.ParseWhere(modelName, opMap)
		if err != nil {
			return
		}
		ql = append(ql, opQ)
	}
	return
}

func (p Parser) ParseWhere(modelName string, data map[string]interface{}) (q Where, err error) {
	m, err := p.tx.Schema().GetInterface(modelName)
	if err != nil {
		return
	}
	q = Where{}
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
		var orQL []Where
		orQL, err = p.parseCompositeQueryList(modelName, orVal)
		if err != nil {
			return
		}
		q.Or = orQL
	}
	if andVal, ok := data["AND"]; ok {
		var andQL []Where
		andQL, err = p.parseCompositeQueryList(modelName, andVal)
		if err != nil {
			return
		}
		q.And = andQL
	}
	if notVal, ok := data["NOT"]; ok {
		var notQL []Where
		notQL, err = p.parseCompositeQueryList(modelName, notVal)
		if err != nil {
			return
		}
		q.Not = notQL
	}
	return
}

func (p Parser) parseSingleRelationshipCriteria(m db.Interface, data map[string]interface{}) (rcl []RelationshipCriterion, err error) {
	rels, err := m.Relationships()
	if err != nil {
		return
	}
	for _, r := range rels {
		if !r.Multi() {
			if value, ok := data[r.Name()]; ok {
				var rc RelationshipCriterion
				rc, err = p.parseRelationshipCriterion(r, value)
				if err != nil {
					return
				}
				rcl = append(rcl, rc)
			}
		}
	}
	return rcl, nil
}

func (p Parser) parseAggregateRelationshipCriteria(m db.Interface, data map[string]interface{}) (arcl []AggregateRelationshipCriterion, err error) {
	rels, err := m.Relationships()
	if err != nil {
		return
	}
	for _, r := range rels {
		if r.Multi() {
			if value, ok := data[r.Name()]; ok {
				var arc AggregateRelationshipCriterion
				arc, err = p.parseAggregateRelationshipCriterion(r, value)
				if err != nil {
					return
				}
				arcl = append(arcl, arc)
			}
		}
	}
	return arcl, nil
}

func parseFieldCriteria(m db.Interface, data map[string]interface{}) (fieldCriteria []FieldCriterion, err error) {
	attrs, err := m.Attributes()
	if err != nil {
		return
	}
	for _, attr := range attrs {
		if value, ok := data[attr.Name()]; ok {
			var fc FieldCriterion
			fc, err = parseFieldCriterion(attr, value)
			fieldCriteria = append(fieldCriteria, fc)
		}
	}
	return
}

func parseFieldCriterion(a db.Attribute, value interface{}) (fc FieldCriterion, err error) {
	fieldName := db.JSONKeyToFieldName(a.Name())

	d := a.Datatype()
	f, err := d.FromJSON()
	if err != nil {
		return
	}

	parsedValue, err := f.Call(value)

	fc = FieldCriterion{
		// TODO handle function values like {startsWith}
		Key: fieldName,
		Val: parsedValue,
	}
	return
}

func (p Parser) parseAggregateRelationshipCriterion(r db.Relationship, value interface{}) (arc AggregateRelationshipCriterion, err error) {
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
		var rc RelationshipCriterion
		rc, err = p.parseRelationshipCriterion(r, v)
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

func (p Parser) parseRelationshipCriterion(r db.Relationship, value interface{}) (rc RelationshipCriterion, err error) {
	mapValue := value.(map[string]interface{})
	m := r.Target()
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
		Relationship: r,
		Where: Where{
			FieldCriteria:                 fc,
			RelationshipCriteria:          rrc,
			AggregateRelationshipCriteria: arrc,
		},
	}
	return
}

func (p Parser) parseInclusion(r db.Relationship, value interface{}) Inclusion {
	if v, ok := value.(bool); ok {
		if v {
			return Inclusion{Relationship: r, Where: Where{}}
		} else {
			panic("Include specified as false?")
		}
	}
	panic("Include with findMany args not yet implemented")
}

func (p Parser) ParseInclude(modelName string, data map[string]interface{}) (i Include, err error) {
	m, err := p.tx.Schema().GetInterface(modelName)
	if err != nil {
		return
	}
	var includes []Inclusion
	rels, err := m.Relationships()
	relsByName := map[string]db.Relationship{}
	for _, r := range rels {
		relsByName[r.Name()] = r
	}

	for k, val := range data {
		r, ok := relsByName[k]
		if !ok {
			err = fmt.Errorf("%w: %v", ErrInvalidRelationship, k)
			return
		}
		inc := p.parseInclusion(r, val)
		includes = append(includes, inc)
	}
	i = Include{Includes: includes}
	return
}
