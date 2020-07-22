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

func (p Parser) consumeData(keys set, data map[string]interface{}) map[string]interface{} {
	var d map[string]interface{}
	if v, ok := data["data"]; ok {
		d = v.(map[string]interface{})
		delete(keys, "data")
	}
	return d
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
