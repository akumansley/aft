package db

type attr struct {
	rec Record
	tx  Tx
}

func (a *attr) ID() ID {
	return a.rec.ID()
}

func (a *attr) Name() string {
	return a.rec.MustGet("name").(string)
}

func (a *attr) Datatype() Datatype {
	ad, _ := a.tx.Schema().GetRelationshipByID(AttributeDatatype.ID)
	dt, _ := a.tx.GetRelatedOne(a.ID(), ad)
	return &coreDatatype{dt, a.tx}
}

func (a *attr) Get(rec Record) interface{} {
	return rec.MustGet(a.Name())
}

func (a *attr) Set(v interface{}, rec Record) {
	f, _ := a.Datatype().FromJSON()
	parsed, _ := f.Call(v)
	rec.Set(a.Name(), parsed)
}
