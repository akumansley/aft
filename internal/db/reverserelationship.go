package db

// Model

var ReverseRelationshipModel = ModelL{
	ID:         MakeID("988e0a48-3da1-40c6-8d1b-be9f44eac5dc"),
	Name:       "relationship",
	Attributes: []AttributeL{},
}
var rrName = ConcreteAttributeL{
	Name:     "name",
	ID:       MakeID("0646d3e7-e8b6-4663-a18c-d9ccf25d61f6"),
	Datatype: String,
}

var ReverseRelationshipReferencing = ConcreteRelationshipL{
	Name:   "referencing",
	ID:     MakeID("fbca6418-da50-4737-ada1-98505dcaec6a"),
	Source: ReverseRelationshipModel,
	Target: ConcreteRelationshipModel,
	Multi:  false,
}

// Loader

type ReverseRelationshipLoader struct{}

func (l ReverseRelationshipLoader) ProvideModel() ModelL {
	return ReverseRelationshipModel
}

func (l ReverseRelationshipLoader) Load(tx Tx, rec Record) Relationship {
	return &reverseRelationship{rec, tx}
}

// Literal

type ReverseRelationshipL struct {
	ID          ID     `record:"id"`
	Name        string `record:"name"`
	Referencing ConcreteRelationshipL
}

func (lit ReverseRelationshipL) GetID() ID {
	return lit.ID
}

func (lit ReverseRelationshipL) MarshalDB() (recs []Record, links []Link) {
	rec := MarshalRecord(lit, ReverseRelationshipModel)
	recs = append(recs, rec)
	links = []Link{
		Link{rec.ID(), lit.Referencing.GetID(), ReverseRelationshipReferencing},
	}
	return
}

// Dynamic

type reverseRelationship struct {
	rec Record
	tx  Tx
}

func (r *reverseRelationship) ID() ID {
	return r.rec.ID()
}

func (r *reverseRelationship) Name() string {
	return rrName.AsAttribute().MustGet(r.rec).(string)
}

func (r *reverseRelationship) Multi() bool {
	refRel, _ := r.tx.Schema().GetRelationshipByID(ReverseRelationshipReferencing.ID)
	referenced, _ := r.tx.GetRelatedOne(r.ID(), refRel)
	return caMulti.AsAttribute().MustGet(referenced).(bool)
}

func (r *reverseRelationship) Source() Interface {
	refRel, _ := r.tx.Schema().GetRelationshipByID(ReverseRelationshipReferencing.ID)
	referenced, _ := r.tx.GetRelatedOne(r.ID(), refRel)

	targetRel, _ := r.tx.Schema().GetRelationshipByID(ConcreteRelationshipTarget.ID)
	mRec, err := r.tx.GetRelatedOne(referenced.ID(), targetRel)
	if err != nil {
		panic("rev source failed")
	}
	return &model{mRec, r.tx}
}

func (r *reverseRelationship) Target() Interface {
	refRel, _ := r.tx.Schema().GetRelationshipByID(ReverseRelationshipReferencing.ID)
	referenced, _ := r.tx.GetRelatedOne(r.ID(), refRel)

	targetRel, _ := r.tx.Schema().GetRelationshipByID(ConcreteRelationshipSource.ID)
	mRec, err := r.tx.GetRelatedOne(referenced.ID(), targetRel)
	if err != nil {
		panic("rev target failed")
	}
	return &model{mRec, r.tx}
}

func (r *reverseRelationship) LoadOne(rec Record) (Record, error) {
	refRel, _ := r.tx.Schema().GetRelationshipByID(ReverseRelationshipReferencing.ID)
	refRec, _ := r.tx.GetRelatedOne(r.ID(), refRel)
	referenced, _ := r.tx.Schema().GetRelationshipByID(refRec.ID())

	if r.Multi() {
		panic("LoadOne on multi record")
	}
	return r.tx.GetRelatedOneReverse(rec.ID(), referenced)
}

func (r *reverseRelationship) LoadMany(rec Record) ([]Record, error) {
	refRel, _ := r.tx.Schema().GetRelationshipByID(ReverseRelationshipReferencing.ID)
	refRec, _ := r.tx.GetRelatedOne(r.ID(), refRel)
	referenced, _ := r.tx.Schema().GetRelationshipByID(refRec.ID())

	if !r.Multi() {
		panic("LoadMany on non-multi record")
	}
	return r.tx.GetRelatedManyReverse(rec.ID(), referenced)
}
