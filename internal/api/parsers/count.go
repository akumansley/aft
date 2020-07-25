package parsers

import (
	"awans.org/aft/internal/api/operations"
	"fmt"
)

func (p Parser) ParseCount(modelName string, args map[string]interface{}) (op operations.CountOperation, err error) {
	m, err := p.Tx.Schema().GetModel(modelName)
	if err != nil {
		return
	}

	unusedKeys := make(set)
	for k := range args {
		unusedKeys[k] = void{}
	}

	where, err := p.consumeWhere(m, unusedKeys, args)
	if err != nil {
		return
	}

	if len(unusedKeys) != 0 {
		return op, fmt.Errorf("%w: %v", ErrUnusedKeys, unusedKeys)
	}
	op = operations.CountOperation{
		Where:   where,
		ModelID: m.ID(),
	}
	return
}
