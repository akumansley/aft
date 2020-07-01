package db

type cdBox struct {
	CoreDatatypeL
}

type CoreDatatypeL struct {
	ID        ID        `record:"id"`
	Name      string    `record:"name"`
	StoredAs  EnumValue `record:"storedAs"`
	Validator Function
}

func (lit CoreDatatypeL) AsDatatype() cdBox {
	return cdBox{lit}
}

func (d cdBox) FromJSON() (Function, error) {
	return d.Validator, nil
}

func (d cdBox) ID() ID {
	return d.CoreDatatypeL.ID
}

func (d cdBox) Name() string {
	return d.CoreDatatypeL.Name
}

func (d cdBox) Storage() EnumValue {
	return d.CoreDatatypeL.StoredAs
}

func (d cdBox) Save(tx RWTx) (rec Record, err error) {
	// TODO
	panic("Not implemnted")
}
