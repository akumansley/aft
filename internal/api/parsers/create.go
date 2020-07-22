package parsers

import (
	"awans.org/aft/internal/api/operations"
	"awans.org/aft/internal/db"
	"fmt"
)

func (p Parser) ParseCreate(modelName string, args map[string]interface{}) (op operations.CreateOperation, err error) {
	m, err := p.Tx.Schema().GetModel(modelName)
	if err != nil {
		return op, fmt.Errorf("%w: %v", ErrInvalidModel, modelName)
	}

	unusedKeys := make(set)
	for k := range args {
		unusedKeys[k] = void{}
	}

	data := p.consumeData(unusedKeys, args)
	rec, nested, err := p.create(m, data)
	if err != nil {
		return
	}

	include, err := p.consumeInclude(modelName, unusedKeys, args)
	if err != nil {
		return
	}

	if len(unusedKeys) != 0 {
		return op, fmt.Errorf("%w: %v", ErrUnusedKeys, unusedKeys)
	}

	return operations.CreateOperation{Record: rec, Nested: nested, Include: include}, nil
}

func (p Parser) parseNestedCreate(rel db.Relationship, data map[string]interface{}) (op operations.NestedOperation, err error) {
	rec, nested, err := p.create(rel.Target(), data)
	if err != nil {
		return
	}
	return operations.NestedCreateOperation{Relationship: rel, Record: rec, Nested: nested}, nil
}

func buildRecordFromData(m db.Interface, keys set, data map[string]interface{}) (db.Record, set, error) {
	rec := db.NewRecord(m)
	attrs, err := m.Attributes()
	if err != nil {
		return nil, keys, err
	}
	for _, a := range attrs {
		ok, err := parseAttribute(a, data, rec)
		if err != nil {
			return nil, nil, err
		}
		if ok {
			delete(keys, a.Name())
		}
	}
	return rec, keys, nil
}

func (p Parser) create(m db.Interface, data map[string]interface{}) (rec db.Record, nested []operations.NestedOperation, err error) {
	unusedKeys := make(set)
	for k := range data {
		unusedKeys[k] = void{}
	}

	rec, unusedKeys, err = buildRecordFromData(m, unusedKeys, data)
	if err != nil {
		return rec, nested, fmt.Errorf("%w: %v", ErrParse, err)
	}

	rels, err := m.Relationships()
	if err != nil {
		return
	}
	nested = []operations.NestedOperation{}
	for _, r := range rels {
		additionalNested, consumed, err := p.parseNestedCreateRelationship(r, data)
		if err != nil {
			return rec, nested, err
		}
		if consumed {
			delete(unusedKeys, r.Name())
		}
		nested = append(nested, additionalNested...)
	}
	if len(unusedKeys) != 0 {
		return rec, nested, fmt.Errorf("%w: %v", ErrUnusedKeys, unusedKeys)
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
