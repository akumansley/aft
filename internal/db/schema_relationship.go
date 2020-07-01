package db

type rel struct {
	rec Record
	tx  Tx
}

func (r *rel) ID() ID {
	return r.rec.ID()
}

func (r *rel) Name() string {
	return r.rec.MustGet("name").(string)
}

func (r *rel) Multi() bool {
	return r.rec.MustGet("multi").(bool)
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
