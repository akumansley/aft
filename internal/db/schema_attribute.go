package db

// attr fronts for any type that implements the Attribute interface
type attr struct {
	rec Record
	tx  Tx
}

func (a *attr) ID() ID {
	return a.rec.ID()
}

func (a *attr) Name() string {
	model := a.tx.Schema().GetModelByID(a.rec.Model().ID())
	nameAttr, err := model.AttributeByName("name")
	if err != nil {
		panic(err)
	}
	return nameAttr.MustGet(a.rec).(string)
}

func (a *attr) Datatype() Datatype {
	ad, _ := a.tx.Schema().GetRelationshipByID(AttributeDatatype.ID)
	dt, _ := a.tx.GetRelatedOne(a.ID(), ad)
	return &coreDatatype{dt, a.tx}
}

func (a *attr) Get(rec Record) (interface{}, error) {
	return rec.get(a.Name())
}

func (a *attr) MustGet(rec Record) interface{} {
	v, err := rec.get(a.Name())
	if err != nil {
		panic(err)
	}
	return v
}

func (a *attr) Set(v interface{}, rec Record) error {
	f, err := a.Datatype().FromJSON()
	if err != nil {
		return err
	}
	parsed, err := f.Call(v)
	if err != nil {
		return err
	}
	rec.set(a.Name(), parsed)
	return nil
}
