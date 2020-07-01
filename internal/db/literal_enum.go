package db

type eBox struct {
	EnumL
}

type EnumL struct {
	ID   ID     `record:"id"`
	Name string `record:"name"`
}

func (lit EnumL) AsDatatype() Datatype {
	return eBox{lit}
}

func (d eBox) FromJSON() Function {
	panic("not implemented")
}

func (d eBox) ID() ID {
	return d.EnumL.ID
}

func (d eBox) Name() string {
	return d.EnumL.Name
}

func (d eBox) Storage() StorageEnumValue {
	return UUIDStorage
}

func (s Schema) SaveEnum(d EnumL) (err error) {
	// TODO
	panic("Not implemnted")
}
