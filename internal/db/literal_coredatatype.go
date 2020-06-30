package db

type cdBox struct {
	CoreDatatype
}

type CoreDatatype struct {
	ID        ID               `record:"id"`
	Name      string           `record:"name"`
	StoredAs  StorageEnumValue `record:"storedAs"`
	Validator Function
}

func (d cdBox) FromJSON() Function {
	return d.Validator
}

func (d cdBox) ID() ID {
	return d.CoreDatatypeL.ID
}

func (d cdBox) Storage() StorageEnumValue {
	return d.CoreDatatypeL.StoredAs
}

func (s Schema) SaveCoreDatatype(d CoreDatatypeL) (err error) {
	// TODO
	panic("Not implemnted")
}
