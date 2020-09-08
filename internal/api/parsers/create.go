package parsers

import (
	"fmt"

	"awans.org/aft/internal/api"
	"awans.org/aft/internal/api/operations"
	"awans.org/aft/internal/db"
)

func (p Parser) ParseCreate(modelName string, args map[string]interface{}) (op operations.CreateOperation, err error) {
	m, err := p.Tx.Schema().GetModel(modelName)
	if err != nil {
		return op, fmt.Errorf("%w: %v", ErrInvalidModel, modelName)
	}

	unusedKeys := make(api.Set)
	for k := range args {
		unusedKeys[k] = api.Void{}
	}

	data := p.consumeData(unusedKeys, args)
	nested, err := p.consumeCreateRel(m, data)
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

	return operations.CreateOperation{
		ModelID: m.ID(),
		Data:    data,
		Nested:  nested,
		FindArgs: operations.FindArgs{
			Include: inc,
			Select:  sel,
		},
	}, nil
}

func (p Parser) parseNestedCreate(rel db.Relationship, data map[string]interface{}) (op operations.NestedOperation, err error) {
	nested, err := p.consumeCreateRel(rel.Target(), data)
	if err != nil {
		return
	}
	return operations.NestedCreateOperation{Relationship: rel, Data: data, Nested: nested}, nil
}

func (p Parser) consumeCreateRel(m db.Interface, data map[string]interface{}) (nested []operations.NestedOperation, err error) {
	unusedKeys := make(api.Set)
	for k := range data {
		unusedKeys[k] = api.Void{}
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
		additionalNested, consumed, err := p.parseNestedCreateRelationship(r, data)
		if err != nil {
			return nested, err
		}
		if consumed {
			delete(unusedKeys, r.Name())
			delete(data, r.Name())
		}
		nested = append(nested, additionalNested...)
	}
	if len(unusedKeys) != 0 {
		return nested, fmt.Errorf("%w: %v", ErrUnusedKeys, unusedKeys)
	}
	return
}

func (p Parser) parseNestedCreateRelationship(r db.Relationship, data map[string]interface{}) ([]operations.NestedOperation, bool, error) {
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
		opList, err := listify(val)

		if err != nil {
			return nil, false, err
		}

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
