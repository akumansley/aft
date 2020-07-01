package db

type evBox struct {
	EnumValueL
}

type EnumValueL struct {
	ID   ID
	Name string
}

func (lit EnumValueL) AsEnumValue() EnumValue {
	return evBox{lit}
}

func (e evBox) ID() ID {
	return e.EnumValueL.ID
}

func (e evBox) Name() string {
	return e.EnumValueL.Name
}
