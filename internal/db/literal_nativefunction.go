package db

type NativeFunctionL struct {
	ID                ID         `record:"id"`
	Name              string     `record:"name"`
	FunctionSignature EnumValueL `record:"functionSignature"`
	Function          Func
}

func (lit NativeFunctionL) GetID() ID {
	return lit.ID
}

func (lit NativeFunctionL) MarshalDB() (recs []Record, links []Link) {
	rec := MarshalRecord(lit, InterfaceModel)
	return
}
