package parsers

import (
	"awans.org/aft/internal/api/operations"
	"awans.org/aft/internal/db"
	"fmt"
)

func (p Parser) ParseFindOne(modelName string, data map[string]interface{}) (op operations.FindOneOperation, err error) {
	m, err := p.Tx.Schema().GetModel(modelName)
	if err != nil {
		return
	}
	var fieldName string
	var value interface{}

	if len(data) > 1 {
		return op, fmt.Errorf("%w: %v", ErrInvalidStructure, data)
	} else if len(data) == 0 {
		return op, fmt.Errorf("%w: %v", ErrInvalidStructure, data)
	}

	for k, v := range data {
		fieldName = db.JSONKeyToFieldName(k)
		var a db.Attribute
		a, err = m.AttributeByName(k)
		if err != nil {
			return
		}
		d := a.Datatype()
		var f db.Function
		f, err = d.FromJSON()
		if err != nil {
			return
		}

		value, err = f.Call(v)
		if err != nil {
			return
		}
	}

	op = operations.FindOneOperation{
		UniqueQuery: operations.UniqueQuery{
			Key: fieldName,
			Val: value,
		},
		ModelID: m.ID(),
	}
	return
}