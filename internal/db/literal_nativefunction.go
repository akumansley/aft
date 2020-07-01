// do i need this file??

package db

type NativeFunctionL struct {
	ID                ID        `record:"id"`
	Name              string    `record:"name"`
	FunctionSignature EnumValue `record:"functionSignature"`
	Function          Func
}

type nBox struct {
	NativeFunctionL
}

func (lit NativeFunctionL) AsFunction() nBox {
	return nBox{lit}
}

func (n nBox) ID() ID {
	return n.NativeFunctionL.ID
}

func (n nBox) Name() string {
	return n.NativeFunctionL.Name
}

func (n nBox) FunctionSignature() EnumValue {
	return n.NativeFunctionL.FunctionSignature
}

func (n nBox) Call(args interface{}) (interface{}, error) {
	return n.NativeFunctionL.Function(args)
}
