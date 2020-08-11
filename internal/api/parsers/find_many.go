package parsers

import (
	"awans.org/aft/internal/api"
	"awans.org/aft/internal/api/operations"
	"awans.org/aft/internal/db"
	"fmt"
)

func (p Parser) ParseFindMany(modelName string, args map[string]interface{}) (op operations.FindManyOperation, err error) {
	m, where, inc, sel, ca, err := p.find(modelName, args)
	if err != nil {
		return
	}
	op = operations.FindManyOperation{
		ModelID: m.ID(),
		FindArgs: operations.FindArgs{
			Where:   where,
			Include: inc,
			Select:  sel,
			Case:    ca,
		},
	}
	return
}

func (p Parser) parseNestedFindMany(modelName string, args map[string]interface{}) (op operations.FindArgs, err error) {
	_, where, inc, sel, ca, err := p.find(modelName, args)
	if err != nil {
		return
	}
	op = operations.FindArgs{
		Where:   where,
		Include: inc,
		Select:  sel,
		Case:    ca,
	}
	return
}

func (p Parser) find(modelName string, args map[string]interface{}) (m db.Interface, where operations.Where, inc operations.Include, sel operations.Select, ca operations.Case, err error) {
	m, err = p.Tx.Schema().GetInterface(modelName)
	if err != nil {
		return
	}

	unusedKeys := make(api.Set)
	for k := range args {
		unusedKeys[k] = api.Void{}
	}

	where, err = p.consumeWhere(m, unusedKeys, args)
	if err != nil {
		return
	}

	inc, sel, err = p.consumeIncludeOrSelect(m, unusedKeys, args)
	if err != nil {
		return
	}

	ca, err = p.consumeCase(m, unusedKeys, args)
	if err != nil {
		return
	}

	if len(unusedKeys) != 0 {
		return m, where, inc, sel, ca, fmt.Errorf("%w: %v", ErrUnusedKeys, unusedKeys)
	}

	return m, where, inc, sel, ca, err
}
