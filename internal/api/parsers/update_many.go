package parsers

import (
	"awans.org/aft/internal/api"
	"awans.org/aft/internal/api/operations"
	"awans.org/aft/internal/db"
	"fmt"
)

func (p Parser) ParseUpdateMany(modelName string, args map[string]interface{}) (op operations.UpdateManyOperation, err error) {
	m, err := p.Tx.Schema().GetModel(modelName)
	if err != nil {
		return
	}

	unusedKeys := make(api.Set)
	for k := range args {
		unusedKeys[k] = api.Void{}
	}

	where, err := p.consumeWhere(m, unusedKeys, args)
	if err != nil {
		return
	}

	data := p.consumeData(unusedKeys, args)
	nested, err := p.consumeUpdateRel(m, data)
	if err != nil {
		return
	}

	if len(unusedKeys) != 0 {
		return op, fmt.Errorf("%w: %v", ErrUnusedKeys, unusedKeys)
	}

	op = operations.UpdateManyOperation{
		ModelID: m.ID(),
		Where:   where,
		Data:    data,
		Nested:  nested,
	}
	return op, err
}

func (p Parser) parseNestedUpdateMany(rel db.Relationship, args map[string]interface{}) (op operations.NestedUpdateManyOperation, err error) {
	unusedKeys := make(api.Set)
	for k := range args {
		unusedKeys[k] = api.Void{}
	}

	where, err := p.consumeWhere(rel.Target(), unusedKeys, args)
	if err != nil {
		return
	}

	data := p.consumeData(unusedKeys, args)
	nested, err := p.consumeUpdateRel(rel.Target(), data)
	if err != nil {
		return
	}

	if len(unusedKeys) != 0 {
		return op, fmt.Errorf("%w: %v", ErrUnusedKeys, unusedKeys)
	}

	op = operations.NestedUpdateManyOperation{Relationship: rel, Data: data, Where: where, Nested: nested}
	return op, err
}
