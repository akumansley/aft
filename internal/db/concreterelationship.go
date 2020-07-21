package db

import (
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

// this breaks a typechecking loop
func init() {
	ConcreteRelationshipModel.Relationships_ = []RelationshipL{
		ConcreteRelationshipTarget,
		ConcreteRelationshipSource,
	}
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

func (l ConcreteRelationshipLoader) Load(tx Tx, rec Record) Relationship {
	return &concreteRelationship{rec, tx}
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

func (lit ConcreteRelationshipL) MarshalDB() (recs []Record, links []Link) {
	rec := MarshalRecord(lit, ConcreteRelationshipModel)
	recs = append(recs, rec)
	target := Link{rec.ID(), lit.Target_.ID(), ConcreteRelationshipTarget}
	links = []Link{target}
	return
}

func (lit ConcreteRelationshipL) ID() ID {
	return lit.ID_
}

func (lit ConcreteRelationshipL) Name() string {
	return lit.Name_
}

func (lit ConcreteRelationshipL) Multi() bool {
	return lit.Multi_
}

func (lit ConcreteRelationshipL) Source() Interface {
	panic("Not implemented")
}

func (lit ConcreteRelationshipL) Target() Interface {
	return lit.Target_
}

func (lit ConcreteRelationshipL) LoadOne(Record) (Record, error) {
	panic("Not implemented")
}

func (lit ConcreteRelationshipL) LoadMany(Record) ([]Record, error) {
	panic("Not implemented")
}

// Dynamic

type concreteRelationship struct {
	rec Record
	tx  Tx
}

func (r *concreteRelationship) ID() ID {
	return r.rec.ID()
}

func (r *concreteRelationship) Name() string {
	return crName.MustGet(r.rec).(string)
}

func (r *concreteRelationship) Multi() bool {
	return crMulti.MustGet(r.rec).(bool)
}

func (r *concreteRelationship) Source() Interface {
	mRec, err := r.tx.getRelatedOneReverse(r.ID(), ModelRelationships.ID())
	if err != nil {
		err = fmt.Errorf("source failed: %w\n", err)
		panic(err)
	}
	return &model{mRec, r.tx}
}

func (r *concreteRelationship) Target() Interface {
	mRec, err := r.tx.getRelatedOne(r.ID(), ConcreteRelationshipTarget.ID())
	if err != nil {
		err = fmt.Errorf("target failed: %v %w\n", r.rec.Map(), err)
		panic(err)
	}
	return &model{mRec, r.tx}
}

func (r *concreteRelationship) LoadOne(rec Record) (Record, error) {
	if r.Multi() {
		panic("LoadOne on multi record")
	}
	return r.tx.getRelatedOne(rec.ID(), r.ID())
}

func (r *concreteRelationship) LoadMany(rec Record) ([]Record, error) {
	if !r.Multi() {
		panic("LoadMany on non-multi record")
	}
	return r.tx.getRelatedMany(rec.ID(), r.ID())
}
