package parsers

import (
	"awans.org/aft/internal/api/operations"
	"awans.org/aft/internal/db"
	"fmt"
)

func (p Parser) ParseUpsert(modelName string, args map[string]interface{}) (op operations.UpsertOperation, err error) {
	m, err := p.Tx.Schema().GetModel(modelName)
	if err != nil {
		return op, fmt.Errorf("%w: %v", ErrInvalidModel, modelName)
	}

	unusedKeys := make(set)
	for k := range args {
		unusedKeys[k] = void{}
	}

	where, err := p.consumeWhere(m, unusedKeys, args)
	if err != nil {
		return
	}

	var create map[string]interface{}
	if v, ok := args["create"]; ok {
		create = v.(map[string]interface{})
		delete(unusedKeys, "create")
	} else {
		return op, fmt.Errorf("%w: missing create", ErrInvalidStructure)
	}
	rec, nestedCreate, err := p.create(m, create)
	if err != nil {
		return
	}

	var update map[string]interface{}
	if v, ok := args["update"]; ok {
		update = v.(map[string]interface{})
		delete(unusedKeys, "update")
	} else {
		return op, fmt.Errorf("%w: missing update", ErrInvalidStructure)
	}
	nestedUpdate, err := p.update(m, update)
	if err != nil {
		return
	}

	include, err := p.consumeInclude(m, unusedKeys, args)
	if err != nil {
		return
	}

	if len(unusedKeys) != 0 {
		return op, fmt.Errorf("%w: %v", ErrUnusedKeys, unusedKeys)
	}

	return operations.UpsertOperation{
		ModelID: m.ID(),
		FindArgs: operations.FindArgs{
			Where:   where,
			Include: include,
		},
		Create:       rec,
		NestedCreate: nestedCreate,
		Update:       update,
		NestedUpdate: nestedUpdate,
	}, nil
}

func (p Parser) parseNestedUpsert(rel db.Relationship, args map[string]interface{}) (op operations.NestedUpsertOperation, err error) {
	unusedKeys := make(set)
	for k := range args {
		unusedKeys[k] = void{}
	}

	where, err := p.consumeWhere(rel.Target(), unusedKeys, args)
	if err != nil {
		return
	}

	var create map[string]interface{}
	if v, ok := args["create"]; ok {
		create = v.(map[string]interface{})
		delete(unusedKeys, "create")
	}
	rec, nestedCreate, err := p.create(rel.Target(), create)
	if err != nil {
		return
	}

	var update map[string]interface{}
	if v, ok := args["update"]; ok {
		create = v.(map[string]interface{})
		delete(unusedKeys, "update")
	}
	nestedUpdate, err := p.update(rel.Target(), update)
	if err != nil {
		return
	}

	if len(unusedKeys) != 0 {
		return op, fmt.Errorf("%w: %v", ErrUnusedKeys, unusedKeys)
	}

	return operations.NestedUpsertOperation{
		Relationship: rel,
		Where:        where,
		Create:       rec,
		NestedCreate: nestedCreate,
		Update:       update,
		NestedUpdate: nestedUpdate,
	}, nil
}
