package db

// Model

var ReverseRelationshipModel = MakeModel(
	MakeID("988e0a48-3da1-40c6-8d1b-be9f44eac5dc"),
	"relationship",
	[]AttributeL{rrName},
	[]RelationshipL{ReverseRelationshipReferencing},
	[]ConcreteInterfaceL{RelationshipInterface},
)

var rrName = MakeConcreteAttribute(
	MakeID("0646d3e7-e8b6-4663-a18c-d9ccf25d61f6"),
	"name",
	String,
)

var ReverseRelationshipReferencing = MakeConcreteRelationship(
	MakeID("fbca6418-da50-4737-ada1-98505dcaec6a"),
	"referencing",
	false,
	ConcreteRelationshipModel,
)

// Loader

type ReverseRelationshipLoader struct{}

func (l ReverseRelationshipLoader) ProvideModel() ModelL {
	return ReverseRelationshipModel
}

func (l ReverseRelationshipLoader) Load(tx Tx, rec Record) Relationship {
	return &reverseRelationship{rec, tx}
}

// Literal

func MakeReverseRelationship(id ID, name string, referencing ConcreteRelationshipL) ReverseRelationshipL {
	return ReverseRelationshipL{
		id, name, referencing,
	}
}

type ReverseRelationshipL struct {
	ID_          ID     `record:"id"`
	Name_        string `record:"name"`
	Referencing_ ConcreteRelationshipL
}

func (lit ReverseRelationshipL) MarshalDB() (recs []Record, links []Link) {
	rec := MarshalRecord(lit, ReverseRelationshipModel)
	recs = append(recs, rec)
	links = []Link{
		Link{rec.ID(), lit.Referencing_.ID(), ReverseRelationshipReferencing},
	}
	return
}

func (lit ReverseRelationshipL) ID() ID {
	return lit.ID_
}

func (lit ReverseRelationshipL) Name() string {
	return lit.Name_
}

func (lit ReverseRelationshipL) Multi() bool {
	return lit.Referencing_.Multi_
}

func (lit ReverseRelationshipL) Source() Interface {
	return lit.Referencing_.Target()
}

func (lit ReverseRelationshipL) Target() Interface {
	return lit.Referencing_.Source()
}

func (lit ReverseRelationshipL) LoadOne(Record) (Record, error) {
	panic("Not implemented")
}

func (lit ReverseRelationshipL) LoadMany(Record) ([]Record, error) {
	panic("Not implemented")
}

func (lit ReverseRelationshipL) LoadOneReverse(Record) (Record, error) {
	panic("Not implemented")
}

func (lit ReverseRelationshipL) LoadManyReverse(Record) ([]Record, error) {
	panic("Not implemented")
}

func (lit ReverseRelationshipL) Connect(Record, Record) error {
	panic("Not implemented")
}

func (lit ReverseRelationshipL) Disconnect(Record, Record) error {
	panic("Not implemented")
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
	return rrName.MustGet(r.rec).(string)
}

func (r *reverseRelationship) Multi() bool {
	// until we have constraints, we assume it's true
	return true
}

func (r *reverseRelationship) Source() Interface {
	referenced, _ := r.tx.getRelatedOne(r.ID(), ReverseRelationshipReferencing.ID())
	mRec, err := r.tx.getRelatedOne(referenced.ID(), ConcreteRelationshipTarget.ID())
	if err != nil {
		panic("rev source failed")
	}
	return &model{mRec, r.tx}
}

func (r *reverseRelationship) Target() Interface {
	referenced, _ := r.tx.getRelatedOne(r.ID(), ReverseRelationshipReferencing.ID())
	mRec, err := r.tx.getRelatedOneReverse(referenced.ID(), ModelRelationships.ID())
	if err != nil {
		panic("rev target failed")
	}
	return &model{mRec, r.tx}
}

func (r *reverseRelationship) LoadOne(rec Record) (Record, error) {
	refRec, _ := r.tx.getRelatedOne(r.ID(), ReverseRelationshipReferencing.ID())
	referenced, _ := r.tx.Schema().GetRelationshipByID(refRec.ID())

	if r.Multi() {
		panic("LoadOne on multi record")
	}
	return r.tx.getRelatedOneReverse(rec.ID(), referenced.ID())
}

func (r *reverseRelationship) LoadMany(rec Record) ([]Record, error) {
	refRec, _ := r.tx.getRelatedOne(r.ID(), ReverseRelationshipReferencing.ID())
	referenced, _ := r.tx.Schema().GetRelationshipByID(refRec.ID())

	if !r.Multi() {
		panic("LoadMany on non-multi record")
	}
	return r.tx.getRelatedManyReverse(rec.ID(), referenced.ID())
}

func (r *reverseRelationship) LoadOneReverse(rec Record) (Record, error) {
	panic("Not implemented")
}

func (r *reverseRelationship) LoadManyReverse(rec Record) ([]Record, error) {
	panic("Not implemented")
}

func (r *reverseRelationship) Connect(c, p Record) error {
	refRec, _ := r.tx.getRelatedOne(r.ID(), ReverseRelationshipReferencing.ID())
	referenced, _ := r.tx.Schema().GetRelationshipByID(refRec.ID())
	return referenced.Connect(p, c)
}

func (r *reverseRelationship) Disconnect(c, p Record) error {
	refRec, _ := r.tx.getRelatedOne(r.ID(), ReverseRelationshipReferencing.ID())
	referenced, _ := r.tx.Schema().GetRelationshipByID(refRec.ID())
	return referenced.Disconnect(p, c)
}
