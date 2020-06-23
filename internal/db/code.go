package db

import (
	"awans.org/aft/internal/datatypes"
)

type Code struct {
	ID                ID
	Name              string
	Runtime           RuntimeEnumValue
	FunctionSignature FunctionSignatureEnumValue
	Code              string
	Function          func(interface{}) (interface{}, error)
	executor          CodeExecutor
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
	uuidValidator.ID:   uuidValidator,
	floatValidator.ID:  floatValidator,
}
