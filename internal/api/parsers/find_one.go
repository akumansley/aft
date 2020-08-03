package parsers

import (
	"awans.org/aft/internal/api"
	"awans.org/aft/internal/api/operations"
	"fmt"
)

func (p Parser) ParseFindOne(modelName string, args map[string]interface{}) (op operations.FindOneOperation, err error) {
	m, err := p.Tx.Schema().GetInterface(modelName)
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

	if len(unusedKeys) != 0 {
		return op, fmt.Errorf("%w: %v", ErrUnusedKeys, unusedKeys)
	}
	op = operations.FindOneOperation{
		ModelID: m.ID(),
		FindArgs: operations.FindArgs{
			Where:   where,
			Include: inc,
			Select:  sel,
		},
	}
	return
}
