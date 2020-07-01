package db

type InterfaceAttributeL struct {
	ID       ID     `record:"id"`
	Name     string `record:"name"`
	Datatype CoreDatatypeL
}

func (lit InterfaceAttributeL) GetID() ID {
	return lit.ID
}

func (lit InterfaceAttributeL) MarshalDB() ([]Record, []Link) {
	rec := MarshalRecord(lit, ConcreteAttributeModel)
	dtl := Link{rec.ID(), lit.Datatype.ID, AttributeDatatype}
	return []Record{rec}, []Link{dtl}
}
