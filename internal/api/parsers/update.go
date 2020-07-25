package parsers

import (
	"awans.org/aft/internal/api/operations"
	"awans.org/aft/internal/db"
	"fmt"
)

func (p Parser) ParseUpdate(modelName string, args map[string]interface{}) (op operations.UpdateOperation, err error) {
	m, err := p.Tx.Schema().GetModel(modelName)
	if err != nil {
		return
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

	data := p.consumeData(unusedKeys, args)
	nested, err := p.consumeUpdateRel(m, data)
	if err != nil {
		return
	}

	if len(unusedKeys) != 0 {
		return op, fmt.Errorf("%w: %v", ErrUnusedKeys, unusedKeys)
	}

	op = operations.UpdateOperation{
		ModelID: m.ID(),
		Where:   where,
		Include: include,
		Data:    data,
		Nested:  nested,
	}
	return op, err
}

func (p Parser) parseNestedUpdate(r db.Relationship, args map[string]interface{}) (op operations.NestedUpdateOperation, err error) {
	unusedKeys := make(set)
	for k := range args {
		unusedKeys[k] = void{}
	}

	where, err := p.consumeWhere(r.Target(), unusedKeys, args)
	if err != nil {
		return
	}

	data := p.consumeData(unusedKeys, args)
	nested, err := p.consumeUpdateRel(r.Target(), data)
	if err != nil {
		return
	}

	op = operations.NestedUpdateOperation{
		Where:        where,
		Relationship: r,
		Data:         data,
		Nested:       nested,
	}
	return op, err
}

func (p Parser) consumeUpdateRel(m db.Interface, data map[string]interface{}) (nested []operations.NestedOperation, err error) {
	unusedKeys := make(set)
	for k := range data {
		unusedKeys[k] = void{}
	}

	// delete all attributes from unusedKeys
	attrs, err := m.Attributes()
	if err != nil {
		return nil, err
	}
	for _, attr := range attrs {
		if _, ok := unusedKeys[attr.Name()]; ok {
			delete(unusedKeys, attr.Name())
		}
	}

	rels, err := m.Relationships()
	if err != nil {
		return
	}
	nested = []operations.NestedOperation{}
	for _, r := range rels {
		additionalNested, consumed, err := p.parseNestedUpdateRelationship(r, data)
		if err != nil {
			return nested, err
		}
		if consumed {
			delete(unusedKeys, r.Name())
			//remove the consumed relationship from the data
			delete(data, r.Name())
		}
		nested = append(nested, additionalNested...)
	}
	if len(unusedKeys) != 0 {
		return nested, fmt.Errorf("%w: %v", ErrUnusedKeys, unusedKeys)
	}
	return
}

func (p Parser) parseNestedUpdateRelationship(r db.Relationship, data map[string]interface{}) ([]operations.NestedOperation, bool, error) {
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
			switch k {
			case "connect":
				nestedOp, ok := op.(map[string]interface{})
				if !ok {
					return nil, false, fmt.Errorf("%w: expected an object, got: %v", ErrInvalidStructure, nestedOp)
				}
				nestedConnect, err := p.parseNestedConnect(r, nestedOp)
				if err != nil {
					return nil, false, err
				}
				nested = append(nested, nestedConnect)
			case "create":
				nestedOp, ok := op.(map[string]interface{})
				if !ok {
					return nil, false, fmt.Errorf("%w: expected an object, got: %v", ErrInvalidStructure, nestedOp)
				}
				nestedCreate, err := p.parseNestedCreate(r, nestedOp)
				if err != nil {
					return nil, false, err
				}
				nested = append(nested, nestedCreate)
			case "delete":
				// different than all the others because relationship : {delete : true}
				//is valid. Doesn't have to be a nested map
				nestedDelete, err := p.parseNestedDelete(r, op)
				if err != nil {
					return nil, false, err
				}
				nested = append(nested, nestedDelete)
			case "deleteMany":
				// different than all the others because relationship : {deleteMany : true}
				//is valid. Doesn't have to be a nested map
				nestedDeleteMany, err := p.parseNestedDeleteMany(r, op)
				if err != nil {
					return nil, false, err
				}
				nested = append(nested, nestedDeleteMany)
			case "update":
				nestedOp, ok := op.(map[string]interface{})
				if !ok {
					return nil, false, fmt.Errorf("%w: expected an object, got: %v", ErrInvalidStructure, nestedOp)
				}
				nestedUpdate, err := p.parseNestedUpdate(r, nestedOp)
				if err != nil {
					return nil, false, err
				}
				nested = append(nested, nestedUpdate)
			case "updateMany":
				nestedOp, ok := op.(map[string]interface{})
				if !ok {
					return nil, false, fmt.Errorf("%w: expected an object, got: %v", ErrInvalidStructure, nestedOp)
				}
				nestedUpdate, err := p.parseNestedUpdateMany(r, nestedOp)
				if err != nil {
					return nil, false, err
				}
				nested = append(nested, nestedUpdate)
			case "upsert":
				nestedOp, ok := op.(map[string]interface{})
				if !ok {
					return nil, false, fmt.Errorf("%w: expected an object, got: %v", ErrInvalidStructure, nestedOp)
				}
				nestedUpsert, err := p.parseNestedUpsert(r, nestedOp)
				if err != nil {
					return nil, false, err
				}
				nested = append(nested, nestedUpsert)
			}
		}
	}

	return nested, true, nil
}
