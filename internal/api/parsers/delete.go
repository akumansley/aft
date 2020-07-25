package parsers

import (
	"awans.org/aft/internal/api/operations"
	"awans.org/aft/internal/db"
	"fmt"
)

func (p Parser) ParseDelete(modelName string, args map[string]interface{}) (op operations.DeleteOperation, err error) {
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

	include, err := p.consumeInclude(m, unusedKeys, args)
	if err != nil {
		return
	}

	nested, err := p.consumeDelete(m, unusedKeys, args)
	if err != nil {
		return
	}

	nested2, err := p.consumeDeleteMany(m, unusedKeys, args)
	if err != nil {
		return
	}

	if len(unusedKeys) != 0 {
		return op, fmt.Errorf("%w: %v", ErrUnusedKeys, unusedKeys)
	}

	return operations.DeleteOperation{
		Where:   where,
		ModelID: m.ID(),
		Include: include,
		Nested:  append(nested, nested2...),
	}, nil

}

func (p Parser) parseNestedDelete(rel db.Relationship, value interface{}) (op operations.NestedDeleteOperation, err error) {
	if v, ok := value.(bool); ok {
		if v {
			return operations.NestedDeleteOperation{Relationship: rel}, nil
		} else {
			return op, fmt.Errorf("%w: delete specified as false", ErrInvalidStructure)
		}
	} else if args, ok := value.(map[string]interface{}); ok {

		unusedKeys := make(set)
		for k := range args {
			unusedKeys[k] = void{}
		}

		where, err := p.consumeWhere(rel.Target(), unusedKeys, args)
		if err != nil {
			return op, err
		}

		nested, err := p.consumeDelete(rel.Target(), unusedKeys, args)
		if err != nil {
			return op, err
		}

		nested2, err := p.consumeDeleteMany(rel.Target(), unusedKeys, args)
		if err != nil {
			return op, err
		}

		if len(unusedKeys) != 0 {
			return op, fmt.Errorf("%w: %v", ErrUnusedKeys, unusedKeys)
		}
		return operations.NestedDeleteOperation{Relationship: rel, Where: where, Nested: append(nested, nested2...)}, nil
	}
	return op, fmt.Errorf("%w: expected an object, got: %v", ErrInvalidStructure, value)

}

func (p Parser) consumeDelete(m db.Interface, keys set, data map[string]interface{}) ([]operations.NestedOperation, error) {
	var w map[string]interface{}
	if v, ok := data["delete"]; ok {
		w = v.(map[string]interface{})
		delete(keys, "delete")
	}
	return p.parseDelete(m, w)
}

func (p Parser) parseDelete(m db.Interface, data map[string]interface{}) (nested []operations.NestedOperation, err error) {
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
		del, err := p.parseNestedDelete(r, val)
		if err != nil {
			return nested, err
		}
		nested = append(nested, del)
	}
	return
}
