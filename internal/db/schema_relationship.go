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

func (r *rel) Source() Model {
	mRec, err := r.tx.GetRelatedOne(r.ID(), RelationshipSource)
	if err != nil {
		panic("source failed")
	}
	return &model{mRec, tx}
}

func (r *rel) Target() Model {
	mRec, err := r.tx.GetRelatedOne(r.ID(), RelationshipTarget)
	if err != nil {
		panic("source failed")
	}
	return &model{mRec, tx}
}
