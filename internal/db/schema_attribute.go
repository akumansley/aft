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

func (a *attr) Getter() Function {
	model := a.tx.Schema().GetModelByID(a.rec.Model().ID())
	getterRel, err := model.RelationshipByName("getter")
	if err != nil {
		panic(err)
	}
	fRec, err := a.tx.GetRelatedOne(a.rec.ID(), getterRel)
	if err != nil {
		panic(err)
	}
	f, err := a.tx.loadFunction(fRec)
	if err != nil {
		panic(err)
	}
	return f
}

func (a *attr) Setter() Function {
	model := a.tx.Schema().GetModelByID(a.rec.Model().ID())
	setterRel, err := model.RelationshipByName("setter")
	if err != nil {
		panic(err)
	}
	fRec, err := a.tx.GetRelatedOne(a.rec.ID(), setterRel)
	if err != nil {
		panic(err)
	}
	f, err := a.tx.loadFunction(fRec)
	if err != nil {
		panic(err)
	}
	return f
}

func (a *attr) Get(rec Record) (interface{}, error) {
	gf := a.Getter()
	args := GetterArgs{rec, a, a.tx}
	val, err := gf.Call(args)
	return val, err
}

func (a *attr) MustGet(rec Record) interface{} {
	v, err := a.Get(rec)
	if err != nil {
		panic(err)
	}
	return v
}

func (a *attr) Set(v interface{}, rec Record) error {
	sf := a.Setter()
	args := SetterArgs{rec, v, a, a.tx}
	_, err := sf.Call(args)
	return err
}
