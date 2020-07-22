package parsers

import (
	"awans.org/aft/internal/api/operations"
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
	Tx db.Tx
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

func (p Parser) parseNestedCreate(rel db.Relationship, data map[string]interface{}) (op operations.NestedOperation, err error) {
	unusedKeys := make(set)
	for k := range data {
		unusedKeys[k] = void{}
	}

	targetModel := rel.Target()
	rec, unusedKeys, err := buildRecordFromData(targetModel, unusedKeys, data)
	if err != nil {
		return
	}
	nested := []operations.NestedOperation{}
	rels, err := targetModel.Relationships()
	if err != nil {
		return
	}
	for _, r := range rels {
		additionalNested, consumed, err := p.parseRelationship(r, data)
		if err != nil {
			return operations.NestedCreateOperation{}, err
		}
		if consumed {
			delete(unusedKeys, r.Name())
		}
		nested = append(nested, additionalNested...)
	}
	if len(unusedKeys) != 0 {
		return operations.NestedCreateOperation{}, fmt.Errorf("%w: %v", ErrUnusedKeys, unusedKeys)
	}
	nestedCreate := operations.NestedCreateOperation{Relationship: rel, Record: rec, Nested: nested}
	return nestedCreate, nil
}

func (p Parser) parseNestedConnect(rel db.Relationship, data map[string]interface{}) (op operations.NestedConnectOperation, err error) {
	if len(data) != 1 {
		panic("Too many keys in a unique query")
	}
	m := rel.Target()

	// this should be a separate method
	var uq operations.UniqueQuery
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
		uq = operations.UniqueQuery{Key: k, Val: val}
	}
	return operations.NestedConnectOperation{Relationship: rel, UniqueQuery: uq}, nil
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

func (p Parser) parseRelationship(r db.Relationship, data map[string]interface{}) ([]operations.NestedOperation, bool, error) {
	// refactor this
	nestedOpMap, ok := data[r.Name()].(map[string]interface{})
	if !ok {
		_, isValue := data[r.Name()]
		if !isValue {
			return []operations.NestedOperation{}, false, nil
		}

		return []operations.NestedOperation{}, false, fmt.Errorf("%w: expected an object, got: %v", ErrInvalidStructure, data)
	}
	var nested []operations.NestedOperation
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
