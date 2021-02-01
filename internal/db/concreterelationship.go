package db

import (
	"encoding/gob"
	"fmt"
)

// Model

var ConcreteRelationshipModel = MakeModel(
	MakeID("90be6901-60a0-4eca-893e-232dc57b0bc1"),
	"concreteRelationship",
	[]AttributeL{
		crName,
		crMulti,
	}, []RelationshipL{},
	[]ConcreteInterfaceL{RelationshipInterface},
)

var ConcreteRelationshipReferencedBy = MakeReverseRelationship(
	MakeID("eade1e6d-c73b-4946-b42d-01c6cff4cf36"),
	"referencedBy",
	ReverseRelationshipReferencing,
)

func init() {
	ConcreteRelationshipModel.Relationships_ = []RelationshipL{
		ConcreteRelationshipTarget,
		ConcreteRelationshipSource,
		ConcreteRelationshipReferencedBy,
	}
	ModelModel.Relationships_ = append(ModelModel.Relationships_, ModelTargeted)
}

var crName = MakeConcreteAttribute(
	MakeID("3e649bba-b5ab-4ee2-a4ef-3da0eed541da"),
	"name",
	String,
)

var crMulti = MakeConcreteAttribute(
	MakeID("3c0b2893-a074-4fd7-931e-9a0e45956b08"),
	"multi",
	Bool,
)

var ConcreteRelationshipTarget = MakeConcreteRelationship(
	MakeID("e194f9bf-ea7a-4c78-a179-bdf9c044ac3c"),
	"target",
	false,
	InterfaceInterface,
)

var ConcreteRelationshipSource = MakeReverseRelationship(
	MakeID("72b8049f-d2ff-4edc-80f6-565dbc1a7d7c"),
	"source",
	ModelRelationships,
)

// Loader

type ConcreteRelationshipLoader struct{}

func (l ConcreteRelationshipLoader) ProvideModel() ModelL {
	return ConcreteRelationshipModel
}

func (l ConcreteRelationshipLoader) Load(rec Record) Relationship {
	return &concreteRelationship{rec}
}

// Literal

// source is determined by the modelL
func MakeConcreteRelationship(id ID, name string, multi bool, target InterfaceL) ConcreteRelationshipL {
	return ConcreteRelationshipL{
		id,
		name,
		multi,
		target,
	}
}

type ConcreteRelationshipL struct {
	ID_     ID     `record:"id"`
	Name_   string `record:"name"`
	Multi_  bool   `record:"multi"`
	Target_ InterfaceL
}

func (lit ConcreteRelationshipL) MarshalDB(b *Builder) (recs []Record, links []Link) {
	rec := MarshalRecord(b, lit)
	recs = append(recs, rec)
	target := Link{rec.ID(), lit.Target_.ID(), ConcreteRelationshipTarget}
	links = []Link{target}
	return
}

func (lit ConcreteRelationshipL) ID() ID {
	return lit.ID_
}

func (lit ConcreteRelationshipL) InterfaceID() ID {
	return ConcreteRelationshipModel.ID()
}

func (lit ConcreteRelationshipL) InterfaceName() string {
	return ConcreteRelationshipModel.Name_
}

func (lit ConcreteRelationshipL) Load(tx Tx) Relationship {
	rel, err := tx.Schema().GetRelationshipByID(lit.ID())
	if err != nil {
		panic(err)
	}
	return rel
}

// Dynamic

func init() {
	gob.Register(&concreteRelationship{})
}

type concreteRelationship struct {
	Rec Record
}

func (r *concreteRelationship) ID() ID {
	return r.Rec.ID()
}

func (r *concreteRelationship) Name() string {
	return r.Rec.MustGet("name").(string)
}

func (r *concreteRelationship) Multi() bool {
	return r.Rec.MustGet("multi").(bool)
}

func (r *concreteRelationship) Source(tx Tx) Interface {
	mRec, err := tx.getRelatedOneReverse(r.ID(), ModelRelationships.ID())
	if err != nil {
		err = fmt.Errorf("source failed: %w\n", err)
		panic(err)
	}
	return &model{mRec}
}

func (r *concreteRelationship) Target(tx Tx) Interface {
	mRec, err := tx.getRelatedOne(r.ID(), ConcreteRelationshipTarget.ID())
	if err != nil {
		err = fmt.Errorf("target failed: %v %w\n", r.Rec.Map(), err)
		panic(err)
	}
	ifc, err := tx.Schema().loadInterface(mRec)
	if err != nil {
		err = fmt.Errorf("target failed: %v %w\n", r.Rec.Map(), err)
		panic(err)
	}
	return ifc
}

func (r *concreteRelationship) LoadOne(tx Tx, rec Record) (Record, error) {
	if r.Multi() {
		panic("LoadOne on multi record")
	}
	return tx.getRelatedOne(rec.ID(), r.ID())
}

func (r *concreteRelationship) LoadMany(tx Tx, rec Record) ([]Record, error) {
	if !r.Multi() {
		panic("LoadMany on non-multi record")
	}
	return tx.getRelatedMany(rec.ID(), r.ID())
}

func (r *concreteRelationship) LoadOneReverse(tx Tx, rec Record) (Record, error) {
	if r.Multi() {
		panic("LoadOne on multi record")
	}
	return tx.getRelatedOneReverse(rec.ID(), r.ID())
}

func (r *concreteRelationship) LoadManyReverse(tx Tx, rec Record) ([]Record, error) {
	if !r.Multi() {
		panic("LoadMany on non-multi record")
	}
	return tx.getRelatedManyReverse(rec.ID(), r.ID())
}

func (r *concreteRelationship) Connect(tx RWTx, p, c Record) error {
	return tx.Connect(p.ID(), c.ID(), r.ID())
}

func (r *concreteRelationship) Disconnect(tx RWTx, p, c Record) error {
	return tx.Disconnect(p.ID(), c.ID(), r.ID())
}
