package parsers

import (
	"awans.org/aft/internal/api/operations"
	"awans.org/aft/internal/db"
	"fmt"
)

func (p Parser) ParseUpdateMany(oldRecs []db.Record, data map[string]interface{}) (op operations.UpdateManyOperation, err error) {
	unusedKeys := make(set)
	for k := range data {
		unusedKeys[k] = void{}
	}

	d := p.consumeData(unusedKeys, data)
	newRecs, err := p.updateMany(oldRecs, d)
	if err != nil {
		return
	}

	if len(unusedKeys) != 0 {
		return op, fmt.Errorf("%w: %v", ErrUnusedKeys, unusedKeys)
	}
	op = operations.UpdateManyOperation{Old: oldRecs, New: newRecs}
	return op, err
}

func (p Parser) parseNestedUpdateMany(oldRecs []db.Record, data map[string]interface{}) (op operations.NestedUpdateManyOperation, err error) {
	newRecs, err := p.updateMany(oldRecs, data)
	if err != nil {
		return
	}
	op = operations.NestedUpdateManyOperation{Old: oldRecs, New: newRecs}
	return op, err
}

func (p Parser) updateMany(oldRecs []db.Record, data map[string]interface{}) (newRecs []db.Record, err error) {
	unusedKeys := make(set)
	for k := range data {
		unusedKeys[k] = void{}
	}
	for _, oldRec := range oldRecs {
		newRec, unusedKeys, err := updateRecordFromData(oldRec, unusedKeys, data)
		if err != nil {
			return newRecs, err
		}
		newRecs = append(newRecs, newRec)
		if len(unusedKeys) != 0 {
			return newRecs, fmt.Errorf("%w: %v", ErrUnusedKeys, unusedKeys)
		}
	}

	return newRecs, err
}
