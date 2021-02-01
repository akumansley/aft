package db

import "fmt"

// Model

var ReverseRelationshipModel = MakeModel(
	MakeID("988e0a48-3da1-40c6-8d1b-be9f44eac5dc"),
	"reverseRelationship",
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

func (l ReverseRelationshipLoader) Load(rec Record) Relationship {
	return &reverseRelationship{rec}
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

func (lit ReverseRelationshipL) MarshalDB(b *Builder) (recs []Record, links []Link) {
	rec := MarshalRecord(b, lit)
	recs = append(recs, rec)
	links = []Link{
		Link{rec.ID(), lit.Referencing_.ID(), ReverseRelationshipReferencing},
	}
	return
}

func (lit ReverseRelationshipL) ID() ID {
	return lit.ID_
}

func (lit ReverseRelationshipL) InterfaceID() ID {
	return ReverseRelationshipModel.ID()
}

func (lit ReverseRelationshipL) InterfaceName() string {
	return ReverseRelationshipModel.Name_
}

func (lit ReverseRelationshipL) Load(tx Tx) Relationship {
	rel, err := tx.Schema().GetRelationshipByID(lit.ID())
	if err != nil {
		panic(err)
	}
	return rel
}

// Dynamic

type reverseRelationship struct {
	rec Record
}

func (r *reverseRelationship) ID() ID {
	return r.rec.ID()
}

func (r *reverseRelationship) Name() string {
	return r.rec.MustGet("name").(string)
}

func (r *reverseRelationship) Multi() bool {
	// until we have constraints, we assume it's true
	return true
}

func (r *reverseRelationship) Source(tx Tx) Interface {
	referenced, err := tx.getRelatedOne(r.ID(), ReverseRelationshipReferencing.ID())
	if err != nil {
		err := fmt.Errorf("rev %v referencing failed \n", r.ID())
		panic(err)
	}
	mRec, err := tx.getRelatedOne(referenced.ID(), ConcreteRelationshipTarget.ID())
	if err != nil {
		panic("rev source failed")
	}
	return &model{mRec}
}

func (r *reverseRelationship) Target(tx Tx) Interface {
	referenced, _ := tx.getRelatedOne(r.ID(), ReverseRelationshipReferencing.ID())
	mRec, err := tx.getRelatedOneReverse(referenced.ID(), ModelRelationships.ID())
	if err != nil {
		panic("rev target failed")
	}
	return &model{mRec}
}

func (r *reverseRelationship) LoadOne(tx Tx, rec Record) (Record, error) {
	refRec, _ := tx.getRelatedOne(r.ID(), ReverseRelationshipReferencing.ID())
	referenced, _ := tx.Schema().GetRelationshipByID(refRec.ID())

	if r.Multi() {
		panic("LoadOne on multi record")
	}
	return tx.getRelatedOneReverse(rec.ID(), referenced.ID())
}

func (r *reverseRelationship) LoadMany(tx Tx, rec Record) ([]Record, error) {
	refRec, _ := tx.getRelatedOne(r.ID(), ReverseRelationshipReferencing.ID())
	referenced, _ := tx.Schema().GetRelationshipByID(refRec.ID())

	if !r.Multi() {
		panic("LoadMany on non-multi record")
	}
	recs, err := tx.getRelatedManyReverse(rec.ID(), referenced.ID())
	return recs, err
}

func (r *reverseRelationship) LoadOneReverse(tx Tx, rec Record) (Record, error) {
	panic("Not implemented")
}

func (r *reverseRelationship) LoadManyReverse(tx Tx, rec Record) ([]Record, error) {
	panic("Not implemented")
}

func (r *reverseRelationship) Connect(tx RWTx, c, p Record) error {
	refRec, _ := tx.getRelatedOne(r.ID(), ReverseRelationshipReferencing.ID())
	referenced, _ := tx.Schema().GetRelationshipByID(refRec.ID())
	return referenced.Connect(tx, p, c)
}

func (r *reverseRelationship) Disconnect(tx RWTx, c, p Record) error {
	refRec, _ := tx.getRelatedOne(r.ID(), ReverseRelationshipReferencing.ID())
	referenced, _ := tx.Schema().GetRelationshipByID(refRec.ID())
	return referenced.Disconnect(tx, p, c)
}
