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
	dt, _ := a.tx.GetRelatedOne(a.ID(), AttributeDatatype)
	return &coreDatatype{dt, a.tx}
}
