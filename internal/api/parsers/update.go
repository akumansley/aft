package parsers

import (
	"awans.org/aft/internal/api/operations"
	"awans.org/aft/internal/db"
	"fmt"
)

func (p Parser) ParseUpdate(oldRec db.Record, args map[string]interface{}) (op operations.UpdateOperation, err error) {
	unusedKeys := make(set)
	for k := range args {
		unusedKeys[k] = void{}
	}

	data := p.consumeData(unusedKeys, args)
	newRec, nested, err := p.update(oldRec, data)
	if err != nil {
		return
	}

	include, err := p.consumeInclude(oldRec.Interface().Name(), unusedKeys, args)
	if err != nil {
		return
	}

	if len(unusedKeys) != 0 {
		return op, fmt.Errorf("%w: %v", ErrUnusedKeys, unusedKeys)
	}
	op = operations.UpdateOperation{Old: oldRec, New: newRec, Nested: nested, Include: include}
	return op, err
}

func updateRecordFromData(oldRec db.Record, keys set, data map[string]interface{}) (db.Record, set, error) {
	newRec := oldRec.DeepCopy()
	attrs, err := oldRec.Interface().Attributes()
	if err != nil {
		return nil, keys, err
	}
	for _, a := range attrs {
		ok, err := parseAttribute(a, data, newRec)
		if err != nil {
			return nil, nil, err
		}
		if ok {
			delete(keys, a.Name())
		}
	}
	return newRec, keys, nil
}

func (p Parser) update(oldRec db.Record, data map[string]interface{}) (newRec db.Record, nested []operations.NestedOperation, err error) {
	unusedKeys := make(set)
	for k := range data {
		unusedKeys[k] = void{}
	}

	newRec, unusedKeys, err = updateRecordFromData(oldRec, unusedKeys, data)

	nested = []operations.NestedOperation{}
	rels, err := oldRec.Interface().Relationships()
	if err != nil {
		return
	}
	for _, r := range rels {
		additionalNested, consumed, err := p.parseNestedUpdateRelationship(r, oldRec, data)
		if err != nil {
			return newRec, nested, err
		}
		if consumed {
			delete(unusedKeys, r.Name())
		}
		nested = append(nested, additionalNested...)
	}

	if len(unusedKeys) != 0 {
		return newRec, nested, fmt.Errorf("%w: %v", ErrUnusedKeys, unusedKeys)
	}
	return
}

func (p Parser) parseNestedUpdateRelationship(r db.Relationship, oldRec db.Record, data map[string]interface{}) ([]operations.NestedOperation, bool, error) {
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
