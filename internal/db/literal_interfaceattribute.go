package db

type iaBox struct {
	InterfaceAttributeL
}

type InterfaceAttributeL struct {
	ID       ID     `record:"id"`
	Name     string `record:"name"`
	Datatype Datatype
}

func (lit InterfaceAttributeL) AsAttribute() Attribute {
	return iaBox{lit}
}

func (a iaBox) ID() ID {
	return a.InterfaceAttributeL.ID
}

func (a iaBox) Name() string {
	return a.InterfaceAttributeL.Name
}

func (a iaBox) Datatype() Datatype {
	return a.InterfaceAttributeL.Datatype
}

func (a iaBox) Get(Record) interface{} {
	panic("Can't get interface attribute")
}

func (a iaBox) Set(interface{}, Record) {
	panic("Can't set interface attribute")
}
