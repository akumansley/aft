package db

type ifBox struct {
	InterfaceL
}

type InterfaceL struct {
	ID         ID     `record:"id"`
	Name       string `record:"name"`
	Attributes []Attribute
}

func (lit InterfaceL) AsInterface() Interface {
	return ifBox{lit}
}

func (m ifBox) ID() ID {
	return m.InterfaceL.ID
}

func (m ifBox) Name() string {
	return m.InterfaceL.Name
}

func (m ifBox) Relationships() ([]Relationship, error) {
	panic("Not implemented")
}

func (m ifBox) Attributes() ([]Attribute, error) {
	return m.InterfaceL.Attributes, nil
}
