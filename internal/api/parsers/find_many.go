package parsers

import (
	"fmt"

	"awans.org/aft/internal/api"
	"awans.org/aft/internal/api/operations"
	"awans.org/aft/internal/db"
)

func (p Parser) ParseFindMany(modelName string, args map[string]interface{}) (op operations.FindManyOperation, err error) {
	m, fa, err := p.find(modelName, args)
	if err != nil {
		return
	}
	op = operations.FindManyOperation{
		ModelID:  m.ID(),
		FindArgs: fa,
	}
	return
}

func (p Parser) parseNestedFindMany(modelName string, args map[string]interface{}) (op operations.FindArgs, err error) {
	_, fa, err := p.find(modelName, args)
	return fa, err
}

func (p Parser) find(modelName string, args map[string]interface{}) (m db.Interface, fa operations.FindArgs, err error) {
	m, err = p.Tx.Schema().GetInterface(modelName)
	if err != nil {
		return
	}

	unusedKeys := make(api.Set)
	for k := range args {
		unusedKeys[k] = api.Void{}
	}

	where, err := p.consumeWhere(m, unusedKeys, args)
	if err != nil {
		return
	}

	inc, sel, err := p.consumeIncludeOrSelect(m, unusedKeys, args)
	if err != nil {
		return
	}

	ca, err := p.consumeCase(m, unusedKeys, args)
	if err != nil {
		return
	}

	fa = operations.FindArgs{
		Where:   where,
		Include: inc,
		Select:  sel,
		Case:    ca,
	}

	if len(unusedKeys) != 0 {
		return m, fa, fmt.Errorf("%w: %v", ErrUnusedKeys, unusedKeys)
	}

	return
}
