package db

type concreteAttr struct {
	rec Record
	tx  Tx
}

func (a *concreteAttr) ID() ID {
	return a.rec.ID()
}

func (a *concreteAttr) Name() string {
	model := a.tx.Schema().GetModelByID(a.rec.Model().ID())
	nameAttr, err := model.AttributeByName("name")
	if err != nil {
		panic(err)
	}
	return nameAttr.MustGet(a.rec).(string)
}

func (a *concreteAttr) Datatype() Datatype {
	ad, _ := a.tx.Schema().GetRelationshipByID(ConcreteAttributeDatatype.ID)
	dt, _ := a.tx.GetRelatedOne(a.ID(), ad)
	return &coreDatatype{dt, a.tx}
}

func (a *concreteAttr) Storage() EnumValue {
	return a.Datatype().Storage()
}

func (a *concreteAttr) Get(rec Record) (interface{}, error) {
	return rec.get(a.Name())
}

func (a *concreteAttr) MustGet(rec Record) interface{} {
	v, err := a.Get(rec)
	if err != nil {
		panic(err)
	}
	return v
}

func (a *concreteAttr) Set(v interface{}, rec Record) error {
	f, err := a.Datatype().FromJSON()
	if err != nil {
		return err
	}
	parsed, err := f.Call(v)
	if err != nil {
		return err
	}
	rec.set(a.Name(), parsed)
	return err
}
