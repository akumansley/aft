package db

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

func (c cBox) Get(Record) (interface{}, error) {
	panic("Not implemented")
}

func (c cBox) MustGet(Record) interface{} {
	panic("Not implemented")
}

func (c cBox) Set(interface{}, Record) error {
	panic("Not implemented")
}
