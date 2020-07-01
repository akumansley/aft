package db

type aBox struct {
	AttributeL
}

type AttributeL struct {
	ID       ID     `record:"id"`
	Name     string `record:"name"`
	Datatype Datatype
}

func (lit AttributeL) AsAttribute() Attribute {
	return aBox{lit}
}

func (a aBox) ID() ID {
	return a.AttributeL.ID
}

func (a aBox) Name() string {
	return a.AttributeL.Name
}

func (a aBox) Datatype() Datatype {
	return a.AttributeL.Datatype
}
