package db

import (
	"awans.org/aft/internal/datatypes"
	"github.com/google/uuid"
)

type Code struct {
	ID                ID
	Name              string
	Runtime           RuntimeEnumValue
	FunctionSignature FunctionSignatureEnumValue
	Code              string
	Function          func(interface{}) (interface{}, error)
	Executor          CodeExecutor
}

type CodeExecutor interface {
	Invoke(Code, interface{}) (interface{}, error)
}

type bootstrapCodeExecutor struct{}

func (*bootstrapCodeExecutor) Invoke(c Code, args interface{}) (interface{}, error) {
	fh := datatypes.GoFunctionHandle{Function: c.Function}
	return fh.Invoke(args)
}

var codeMap map[ID]Code = map[ID]Code{
	boolValidator.ID:   boolValidator,
	intValidator.ID:    intValidator,
	stringValidator.ID: stringValidator,
	textValidator.ID:   textValidator,
	uuidValidator.ID:   uuidValidator,
	floatValidator.ID:  floatValidator,
}

func SaveCode(storeCode Record, c Code) error {
	ew := NewRecordWriter(storeCode)
	ew.Set("id", uuid.UUID(c.ID))
	ew.Set("name", c.Name)
	ew.Set("runtime", uuid.UUID(c.Runtime.ID))
	ew.Set("functionSignature", uuid.UUID(c.FunctionSignature.ID))
	ew.Set("code", c.Code)
	return ew.err
}

func RecordToCode(r Record, tx Tx) (Code, error) {
	rt, err := RecordToEnumValue(r, "runtime", tx)
	if err != nil {
		return Code{}, err
	}
	fs, err := RecordToEnumValue(r, "functionSignature", tx)
	if err != nil {
		return Code{}, err
	}
	ew := NewRecordWriter(r)
	c := Code{
		ID:                r.ID(),
		Name:              ew.Get("name").(string),
		Runtime:           RuntimeEnumValue{rt},
		Code:              ew.Get("code").(string),
		FunctionSignature: FunctionSignatureEnumValue{fs},
		Executor:          tx.Ex(),
	}
	if ew.err != nil {
		return Code{}, err
	}
	if _, ok := codeMap[c.ID]; ok {
		c.Function = codeMap[c.ID].Function
	}
	return c, nil
}
