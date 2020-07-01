package db

type CoreDatatypeL struct {
	ID        ID         `record:"id"`
	Name      string     `record:"name"`
	StoredAs  EnumValueL `record:"storedAs"`
	Validator NativeFunctionL
}

func (lit CoreDatatypeL) GetID() ID {
	return lit.ID
}

func (lit CoreDatatypeL) MarshalDB() ([]Record, []Link) {
	rec := MarshalRecord(lit, CoreDatatypeModel)
	dtl := Link{rec.ID(), lit.Validator.ID, DatatypeValidator}
	return []Record{rec}, []Link{dtl}
}

func (lit CoreDatatypeL) AsDatatype() Datatype {
	return cdBox{lit}
}

type cdBox struct {
	CoreDatatypeL
}

func (c cdBox) ID() ID {
	return c.CoreDatatypeL.ID
}
func (c cdBox) Name() string {
	return c.CoreDatatypeL.Name
}

func (c cdBox) Storage() EnumValue {
	return c.StoredAs.AsEnumValue()

}

func (c cdBox) FromJSON() (Function, error) {
	panic("Not implemented")
}
