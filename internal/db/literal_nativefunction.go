package db

type NativeFunctionL struct {
	ID                ID                         `record:"id"`
	Name              string                     `record:"name"`
	FunctionSignature FunctionSignatureEnumValue `record:"functionSignature"`
	Function          func(interface{}) (interface{}, error)
}

// Runtime:           Native,
