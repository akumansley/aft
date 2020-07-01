package db

type EnumValueL struct {
	ID   ID
	Name string
}

func (lit EnumValueL) GetID() ID {
	return lit.ID
}

func (lit EnumValueL) MarshalDB() (recs []Record, links []Link) {
	rec := MarshalRecord(lit, EnumValueModel)
	return []Record{rec}, []Link{}
}

func (lit EnumValueL) AsEnumValue() EnumValue {
	return evBox{lit}
}

type evBox struct {
	EnumValueL
}

func (e evBox) ID() ID {
	return e.EnumValueL.ID
}

func (e evBox) Name() string {
	return e.EnumValueL.Name
}
