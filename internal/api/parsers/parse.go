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

func (p Parser) consumeData(keys set, data map[string]interface{}) map[string]interface{} {
	var d map[string]interface{}
	if v, ok := data["data"]; ok {
		d = v.(map[string]interface{})
		delete(keys, "data")
	}
	return d
}

func (p Parser) parseNestedConnect(rel db.Relationship, data map[string]interface{}) (op operations.NestedConnectOperation, err error) {
	where, err := p.parseWhere(rel.Target(), data)
	if err != nil {
		return op, err
	}
	return operations.NestedConnectOperation{Relationship: rel, Where: where}, nil
}

func (p Parser) parseNestedDisconnect(rel db.Relationship, data map[string]interface{}) (op operations.NestedDisconnectOperation, err error) {
	where, err := p.parseWhere(rel.Target(), data)
	if err != nil {
		return op, err
	}
	return operations.NestedDisconnectOperation{Relationship: rel, Where: where}, nil
}

func listify(val interface{}) []interface{} {
	var opList []interface{}
	switch v := val.(type) {
	case map[string]interface{}:
		opList = []interface{}{v}
	case []interface{}:
		opList = v
	case interface{}:
		opList = []interface{}{v}
	default:
		panic("Invalid input")
	}
	return opList
}
