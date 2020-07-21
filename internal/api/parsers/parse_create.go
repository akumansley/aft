package parsers

import (
	"awans.org/aft/internal/api/operations"
	"awans.org/aft/internal/db"
	"fmt"
)

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

func (p Parser) ParseCreate(modelName string, data map[string]interface{}) (op operations.CreateOperation, err error) {
	unusedKeys := make(set)
	for k := range data {
		unusedKeys[k] = void{}
	}

	m, err := p.Tx.Schema().GetModel(modelName)
	if err != nil {
		return op, fmt.Errorf("%w: %v", ErrInvalidModel, modelName)
	}
	rec, unusedKeys, err := buildRecordFromData(m, unusedKeys, data)
	if err != nil {
		return op, fmt.Errorf("%w: %v", ErrParse, err)
	}
	nested := []operations.NestedOperation{}
	rels, err := m.Relationships()
	if err != nil {
		return
	}
	for _, r := range rels {
		additionalNested, consumed, err := p.parseRelationship(r, data)
		if err != nil {
			return op, err
		}
		if consumed {
			delete(unusedKeys, r.Name())
		}
		nested = append(nested, additionalNested...)
	}
	if len(unusedKeys) != 0 {
		return op, fmt.Errorf("%w: %v", ErrUnusedKeys, unusedKeys)
	}
	op = operations.CreateOperation{Record: rec, Nested: nested}
	return op, err
}