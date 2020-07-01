package db

type aBox struct {
	ConcreteAttributeL
}

type ConcreteAttributeL struct {
	ID       ID     `record:"id"`
	Name     string `record:"name"`
	Datatype Datatype
}

func (lit ConcreteAttributeL) AsAttribute() Attribute {
	return aBox{lit}
}

func (a aBox) ID() ID {
	return a.ConcreteAttributeL.ID
}

func (a aBox) Name() string {
	return a.ConcreteAttributeL.Name
}

func (a aBox) Datatype() Datatype {
	return a.ConcreteAttributeL.Datatype
}
