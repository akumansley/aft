package db

type concreteRelationship struct {
	rec Record
	tx  Tx
}

func (r *concreteRelationship) ID() ID {
	return r.rec.ID()
}

func (r *concreteRelationship) Name() string {
	model := r.tx.Schema().GetModelByID(r.rec.Model().ID())
	nameAttr, err := model.AttributeByName("name")
	if err != nil {
		panic(err)
	}
	return nameAttr.MustGet(r.rec).(string)
}

func (r *concreteRelationship) Multi() bool {
	model := r.tx.Schema().GetModelByID(r.rec.Model().ID())
	multiAttr, err := model.AttributeByName("multi")
	if err != nil {
		panic(err)
	}
	return multiAttr.MustGet(r.rec).(bool)
}

func (r *concreteRelationship) Source() Interface {
	sourceRel, _ := r.tx.Schema().GetRelationshipByID(ConcreteRelationshipSource.ID)
	mRec, err := r.tx.GetRelatedOne(r.ID(), sourceRel)
	if err != nil {
		panic("source failed")
	}
	return &model{mRec, r.tx}
}

func (r *concreteRelationship) Target() Interface {
	targetRel, _ := r.tx.Schema().GetRelationshipByID(ConcreteRelationshipTarget.ID)
	mRec, err := r.tx.GetRelatedOne(r.ID(), targetRel)
	if err != nil {
		panic("source failed")
	}
	return &model{mRec, r.tx}
}

func (r *concreteRelationship) LoadOne(rec Record) (Record, error) {
	if r.Multi() {
		panic("LoadOne on multi record")
	}
	return r.tx.GetRelatedOne(rec.ID(), r)
}

func (r *concreteRelationship) LoadMany(rec Record) ([]Record, error) {
	if !r.Multi() {
		panic("LoadMany on non-multi record")
	}
	return r.tx.GetRelatedMany(rec.ID(), r)
}
