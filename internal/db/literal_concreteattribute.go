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

func (a aBox) Get(rec Record) interface{} {
	return rec.MustGet(a.Name())
}

func (a aBox) Set(v interface{}, rec Record) {
	f, _ := a.Datatype().FromJSON()
	parsed, _ := f.Call(v)
	rec.Set(a.Name(), parsed)
}
