package db

type nBox struct {
	NativeFunctionL
}

type NativeFunctionL struct {
	ID                ID                         `record:"id"`
	Name              string                     `record:"name"`
	FunctionSignature FunctionSignatureEnumValue `record:"functionSignature"`
	Function          func(interface{}) (interface{}, error)
}

func (lit NativeFunctionL) AsFunction() Function {
	return nBox{lit}
}

func (n nBox) ID() ID {
	return n.NativeFunctionL.ID
}

func (n nBox) Name() string {
	return n.NativeFunctionL.Name
}

func (n nBox) Runtime() RuntimeEnumValue {
	return Native
}

func (n nBox) FunctionSignature() FunctionSignatureEnumValue {
	return n.NativeFunctionL.FunctionSignature
}

func (n nBox) Call(args interface{}) (interface{}, error) {
	return n.NativeFunctionL.Function(args)
}
