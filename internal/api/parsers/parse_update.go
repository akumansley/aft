package parsers

import (
	"awans.org/aft/internal/api/operations"
	"awans.org/aft/internal/db"
	"fmt"
)

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

func (p Parser) ParseUpdate(oldRec db.Record, data map[string]interface{}) (op operations.UpdateOperation, err error) {
	unusedKeys := make(set)
	for k := range data {
		unusedKeys[k] = void{}
	}

	newRec, unusedKeys, err := updateRecordFromData(oldRec, unusedKeys, data)
	if len(unusedKeys) != 0 {
		return op, fmt.Errorf("%w: %v", ErrUnusedKeys, unusedKeys)
	}
	op = operations.UpdateOperation{Old: oldRec, New: newRec}
	return op, err
}

func (p Parser) ParseUpdateMany(oldRecs []db.Record, data map[string]interface{}) (op operations.UpdateManyOperation, err error) {
	unusedKeys := make(set)
	for k := range data {
		unusedKeys[k] = void{}
	}
	var newRecs []db.Record
	for _, oldRec := range oldRecs {
		newRec, unusedKeys, err := updateRecordFromData(oldRec, unusedKeys, data)
		if err != nil {
			return op, err
		}
		newRecs = append(newRecs, newRec)
		if len(unusedKeys) != 0 {
			return op, fmt.Errorf("%w: %v", ErrUnusedKeys, unusedKeys)
		}
	}
	op = operations.UpdateManyOperation{Old: oldRecs, New: newRecs}
	return op, err
}