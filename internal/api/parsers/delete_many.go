package parsers

import (
	"awans.org/aft/internal/api/operations"
	"awans.org/aft/internal/db"
	"fmt"
)

func (p Parser) ParseDeleteMany(modelName string, args map[string]interface{}) (op operations.DeleteManyOperation, err error) {
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

	op = operations.DeleteManyOperation{
		Where:   where,
		ModelID: m.ID(),
		Nested:  append(nested, nested2...),
	}
	return
}

func (p Parser) parseNestedDeleteMany(r db.Relationship, value interface{}) (op operations.NestedDeleteManyOperation, err error) {
	if v, ok := value.(bool); ok {
		if v {
			return operations.NestedDeleteManyOperation{Relationship: r}, nil
		} else {
			return op, fmt.Errorf("%w: delete specified as false", ErrInvalidStructure)
		}
	} else if args, ok := value.(map[string]interface{}); ok {
		unusedKeys := make(set)
		for k := range args {
			unusedKeys[k] = void{}
		}

		where, err := p.consumeWhere(r.Target(), unusedKeys, args)
		if err != nil {
			return op, err
		}

		nested, err := p.consumeDeleteMany(r.Target(), unusedKeys, args)
		if err != nil {
			return op, err
		}

		nested2, err := p.consumeDelete(r.Target(), unusedKeys, args)
		if err != nil {
			return op, err
		}

		if len(unusedKeys) != 0 {
			return op, fmt.Errorf("%w: %v", ErrUnusedKeys, unusedKeys)
		}

		op = operations.NestedDeleteManyOperation{
			Where:        where,
			Relationship: r,
			Nested:       append(nested, nested2...),
		}
		return op, err
	}
	return op, fmt.Errorf("%w: expected an object, got: %v", ErrInvalidStructure, value)
}

func (p Parser) consumeDeleteMany(m db.Interface, keys set, data map[string]interface{}) ([]operations.NestedOperation, error) {
	var w map[string]interface{}
	if v, ok := data["deleteMany"]; ok {
		w = v.(map[string]interface{})
		delete(keys, "deleteMany")
	}
	return p.parseDeleteMany(m, w)
}

func (p Parser) parseDeleteMany(m db.Interface, data map[string]interface{}) (nested []operations.NestedOperation, err error) {
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
		del, err := p.parseNestedDeleteMany(r, val)
		if err != nil {
			return nested, err
		}
		nested = append(nested, del)
	}
	return
}
