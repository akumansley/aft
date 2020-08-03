package parsers

import (
	"awans.org/aft/internal/api"
	"awans.org/aft/internal/api/operations"
	"awans.org/aft/internal/db"
	"fmt"
)

func (p Parser) ParseUpsert(modelName string, args map[string]interface{}) (op operations.UpsertOperation, err error) {
	m, err := p.Tx.Schema().GetModel(modelName)
	if err != nil {
		return op, fmt.Errorf("%w: %v", ErrInvalidModel, modelName)
	}

	unusedKeys := make(api.Set)
	for k := range args {
		unusedKeys[k] = api.Void{}
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
	nestedCreate, err := p.consumeCreateRel(m, create)
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
	nestedUpdate, err := p.consumeUpdateRel(m, update)
	if err != nil {
		return
	}

	inc, sel, err := p.consumeIncludeOrSelect(m, unusedKeys, args)
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
			Include: inc,
			Select:  sel,
		},
		Create:       create,
		NestedCreate: nestedCreate,
		Update:       update,
		NestedUpdate: nestedUpdate,
	}, nil
}

func (p Parser) parseNestedUpsert(rel db.Relationship, args map[string]interface{}) (op operations.NestedUpsertOperation, err error) {
	unusedKeys := make(api.Set)
	for k := range args {
		unusedKeys[k] = api.Void{}
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
	nestedCreate, err := p.consumeCreateRel(rel.Target(), create)
	if err != nil {
		return
	}

	var update map[string]interface{}
	if v, ok := args["update"]; ok {
		create = v.(map[string]interface{})
		delete(unusedKeys, "update")
	}
	nestedUpdate, err := p.consumeUpdateRel(rel.Target(), update)
	if err != nil {
		return
	}

	if len(unusedKeys) != 0 {
		return op, fmt.Errorf("%w: %v", ErrUnusedKeys, unusedKeys)
	}

	return operations.NestedUpsertOperation{
		Relationship: rel,
		Where:        where,
		Create:       create,
		NestedCreate: nestedCreate,
		Update:       update,
		NestedUpdate: nestedUpdate,
	}, nil
}
