package db

type rel struct {
	rec Record
	tx  Tx
}

func (r *rel) ID() ID {
	return r.rec.ID()
}

func (r *rel) Name() string {
	model := r.tx.Schema().GetModelByID(r.rec.Model().ID())
	nameAttr, err := model.AttributeByName("name")
	if err != nil {
		panic(err)
	}
	return nameAttr.MustGet(r.rec).(string)
}

func (r *rel) Multi() bool {
	model := r.tx.Schema().GetModelByID(r.rec.Model().ID())
	multiAttr, err := model.AttributeByName("multi")
	if err != nil {
		panic(err)
	}
	return multiAttr.MustGet(r.rec).(bool)
}

func (r *rel) Source() Interface {
	sourceRel, _ := r.tx.Schema().GetRelationshipByID(RelationshipSource.ID)
	mRec, err := r.tx.GetRelatedOne(r.ID(), sourceRel)
	if err != nil {
		panic("source failed")
	}
	return &model{mRec, r.tx}
}

func (r *rel) Target() Interface {
	targetRel, _ := r.tx.Schema().GetRelationshipByID(RelationshipTarget.ID)
	mRec, err := r.tx.GetRelatedOne(r.ID(), targetRel)
	if err != nil {
		panic("source failed")
	}
	return &model{mRec, r.tx}
}
