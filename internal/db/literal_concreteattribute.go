package db

var ConcreteAttributeModel = ModelL{
	ID:   MakeID("14d840f5-344f-4e23-af12-d4caa1ffa848"),
	Name: "concreteAttribute",
	Attributes: []AttributeL{
		ConcreteAttributeL{
			Name:     "name",
			ID:       MakeID("51605ada-5326-4cfd-9f31-f10bc4dfbf03"),
			Datatype: String,
		},
	},
}

type GetterArgs struct {
	rec  Record
	attr Attribute
	tx   Tx
}

func getter(value interface{}) (interface{}, error) {
	args := value.(GetterArgs)
	rec := args.rec
	attr := args.attr
	return rec.get(attr.Name())
}

var concreteGetter = NativeFunctionL{
	Name:              "concreteGetter",
	ID:                MakeID("532e86b2-0a9f-498b-abc2-0005dd6c8d71"),
	Function:          getter,
	FunctionSignature: Getter,
}

type SetterArgs struct {
	rec   Record
	value interface{}
	attr  Attribute
	tx    Tx
}

func setter(value interface{}) (interface{}, error) {
	args := value.(SetterArgs)
	a := args.attr
	rec := args.rec
	v := args.value

	f, err := a.Datatype().FromJSON()
	if err != nil {
		return nil, err
	}
	parsed, err := f.Call(v)
	if err != nil {
		return nil, err
	}
	rec.set(a.Name(), parsed)
	return nil, err
}

var concreteSetter = NativeFunctionL{
	Name:              "concreteSetter",
	ID:                MakeID("ab501a81-7e8f-4c3e-bc64-b2f3586da8ed"),
	Function:          setter,
	FunctionSignature: Setter,
}

type ConcreteAttributeL struct {
	ID       ID     `record:"id"`
	Name     string `record:"name"`
	Datatype DatatypeL
}

func (lit ConcreteAttributeL) GetID() ID {
	return lit.ID
}

func (lit ConcreteAttributeL) MarshalDB() ([]Record, []Link) {
	rec := MarshalRecord(lit, ConcreteAttributeModel)
	dtl := Link{rec.ID(), lit.Datatype.GetID(), AttributeDatatype}
	return []Record{rec}, []Link{dtl}
}

func (lit ConcreteAttributeL) AsAttribute() Attribute {
	return cBox{lit}
}

type cBox struct {
	ConcreteAttributeL
}

func (c cBox) ID() ID {
	return c.ConcreteAttributeL.ID
}

func (c cBox) Name() string {
	return c.ConcreteAttributeL.Name
}

func (c cBox) Datatype() Datatype {
	return c.ConcreteAttributeL.Datatype.AsDatatype()
}

func (c cBox) Getter() Function {
	panic("Not implemented")
}

func (c cBox) Setter() Function {
	panic("Not implemented")
}

func (c cBox) Get(Record) (interface{}, error) {
	panic("Not implemented")
}

func (c cBox) MustGet(Record) interface{} {
	panic("Not implemented")
}

func (c cBox) Set(interface{}, Record) error {
	panic("Not implemented")
}
