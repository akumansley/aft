package parsers

import (
	"awans.org/aft/internal/api/operations"
	"awans.org/aft/internal/db"
	"fmt"
)

func (p Parser) ParseFindMany(modelName string, args map[string]interface{}) (op operations.FindManyOperation, err error) {
	m, where, include, err := p.find(modelName, args)
	if err != nil {
		return
	}
	op = operations.FindManyOperation{
		Where:   where,
		ModelID: m.ID(),
		Include: include,
	}
	return
}

func (p Parser) parseNestedFindMany(modelName string, args map[string]interface{}) (op operations.NestedFindManyOperation, err error) {
	_, where, include, err := p.find(modelName, args)
	if err != nil {
		return
	}
	op = operations.NestedFindManyOperation{
		Where:   where,
		Include: include,
	}
	return
}

func (p Parser) find(modelName string, args map[string]interface{}) (m db.Model, where operations.Where, include operations.Include, err error) {
	m, err = p.Tx.Schema().GetModel(modelName)
	if err != nil {
		return
	}

	unusedKeys := make(set)
	for k := range args {
		unusedKeys[k] = void{}
	}

	where, err = p.consumeWhere(modelName, unusedKeys, args)
	if err != nil {
		return
	}

	include, err = p.consumeInclude(modelName, unusedKeys, args)
	if err != nil {
		return
	}

	if len(unusedKeys) != 0 {
		return m, where, include, fmt.Errorf("%w: %v", ErrUnusedKeys, unusedKeys)
	}

	return m, where, include, err
}
