package db

import (
	"encoding/gob"
	"errors"
	"fmt"
)

// Model

var InterfaceRelationshipModel = MakeModel(
	MakeID("15753e7e-9ea1-4e22-94c9-ca5dafcb1e5d"),
	"interfaceRelationship",
	[]AttributeL{
		irName, irMulti,
	},
	[]RelationshipL{},
	[]ConcreteInterfaceL{RelationshipInterface},
)

func init() {
	InterfaceRelationshipModel.Relationships_ = []RelationshipL{
		InterfaceRelationshipTarget,
		InterfaceRelationshipSource,
	}
}

var irName = MakeConcreteAttribute(
	MakeID("04cacda4-5121-4343-92e2-30c1fbf7dec8"),
	"name",
	String,
)

var irMulti = MakeConcreteAttribute(
	MakeID("72787dca-8845-42d6-83ce-807d2f1ae899"),
	"multi",
	Bool,
)

var InterfaceRelationshipTarget = MakeConcreteRelationship(
	MakeID("4aec46de-4061-492a-ae72-6621aa120b39"),
	"target",
	false,
	InterfaceInterface,
)

var InterfaceRelationshipSource = MakeReverseRelationship(
	MakeID("40c86d89-5282-476c-a827-50591229e414"),
	"source",
	InterfaceRelationships,
)

// Loader
type InterfaceRelationshipLoader struct{}

func (l InterfaceRelationshipLoader) ProvideModel() ModelL {
	return InterfaceRelationshipModel
}

func (l InterfaceRelationshipLoader) Load(rec Record) Relationship {
	return &interfaceRelationship{rec}
}

// Literal

// source is determined by the modelL
func MakeInterfaceRelationship(id ID, name string, multi bool, target InterfaceL) InterfaceRelationshipL {
	return InterfaceRelationshipL{
		id,
		name,
		multi,
		nil,
		target,
	}
}

func MakeInterfaceRelationshipWithSource(id ID, name string, multi bool, source, target InterfaceL) InterfaceRelationshipL {
	return InterfaceRelationshipL{
		id,
		name,
		multi,
		source,
		target,
	}
}

type InterfaceRelationshipL struct {
	ID_     ID     `record:"id"`
	Name_   string `record:"name"`
	Multi_  bool   `record:"multi"`
	Source_ InterfaceL
	Target_ InterfaceL
}

func (lit InterfaceRelationshipL) MarshalDB(b *Builder) (recs []Record, links []Link) {
	rec := MarshalRecord(b, lit)
	recs = append(recs, rec)
	target := Link{lit, lit.Target_, InterfaceRelationshipTarget}
	links = []Link{target}
	if lit.Source_ != nil {
		source := Link{lit.Source_, lit, InterfaceRelationships}
		links = append(links, source)
	}
	return
}

func (lit InterfaceRelationshipL) ID() ID {
	return lit.ID_
}

func (lit InterfaceRelationshipL) InterfaceID() ID {
	return InterfaceRelationshipModel.ID()
}

func (lit InterfaceRelationshipL) InterfaceName() string {
	return InterfaceRelationshipModel.Name_
}

func (lit InterfaceRelationshipL) Load(tx Tx) Relationship {
	rel, err := tx.Schema().GetRelationshipByID(lit.ID())
	if err != nil {
		panic(err)
	}
	return rel
}

// Dynamic

func init() {
	gob.Register(&interfaceRelationship{})
}

type interfaceRelationship struct {
	Rec Record
}

func (r *interfaceRelationship) ID() ID {
	return r.Rec.ID()
}

func (r *interfaceRelationship) Name() string {
	return r.Rec.MustGet("name").(string)
}

func (r *interfaceRelationship) Multi() bool {
	return r.Rec.MustGet("multi").(bool)
}

func (r *interfaceRelationship) Source(tx Tx) Interface {
	mRec, err := tx.getRelatedOneReverse(r.ID(), InterfaceRelationships.ID())
	if err != nil {
		err = fmt.Errorf("source failed: %w\n", err)
		panic(err)
	}
	return &iface{mRec}
}

func (r *interfaceRelationship) Target(tx Tx) Interface {
	mRec, err := tx.getRelatedOne(r.ID(), InterfaceRelationshipTarget.ID())
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

func (r *interfaceRelationship) getImplementingRelationship(tx Tx, interfaceID ID) (rel Relationship, err error) {
	iface, err := tx.Schema().GetInterfaceByID(interfaceID)
	if err != nil {
		return
	}
	rel, err = iface.RelationshipByName(tx, r.Name())
	if errors.Is(err, ErrNotFound) {
		err = fmt.Errorf("%w: %v does not implement %v - no relationship %v\n", err, iface.Name(), r.Source(tx).Name(), r.Name())
	}
	return
}

func (r *interfaceRelationship) getAllImplementingRelationships(tx Tx) (rels []Relationship, err error) {
	sourceIface := r.Source(tx)
	models := tx.Ref(ModelModel.ID())
	interfaces := tx.Ref(InterfaceModel.ID())
	mRecs := tx.Query(models,
		Join(interfaces, models.Rel(ModelImplements.Load(tx))),
		Filter(interfaces, EqID(sourceIface.ID())),
	).Records()
	for _, rec := range mRecs {
		m := tx.Schema().LoadModel(rec)
		rel, err := m.RelationshipByName(tx, r.Name())
		if err != nil {
			if errors.Is(err, ErrNotFound) {
				err = fmt.Errorf("%w: %v does not implement %v - no relationship %v\n", err, m.Name(), sourceIface.Name(), r.Name())
			}
			return nil, err
		}
		rels = append(rels, rel)
	}

	return
}

func (r *interfaceRelationship) LoadOne(tx Tx, rec Record) (Record, error) {
	if r.Multi() {
		panic("LoadOne on multi record")
	}

	rel, err := r.getImplementingRelationship(tx, rec.InterfaceID())
	if err != nil {
		return nil, err
	}
	return rel.LoadOne(tx, rec)
}

func (r *interfaceRelationship) LoadMany(tx Tx, rec Record) ([]Record, error) {
	if !r.Multi() {
		panic("LoadMany on non-multi record")
	}
	rel, err := r.getImplementingRelationship(tx, rec.InterfaceID())
	if err != nil {
		return nil, err
	}
	return rel.LoadMany(tx, rec)
}

func (r *interfaceRelationship) LoadManyReverse(tx Tx, rec Record) (recs []Record, err error) {
	rels, err := r.getAllImplementingRelationships(tx)
	if err != nil {
		return nil, err
	}
	for _, rel := range rels {
		var newRecs []Record
		newRecs, err = rel.LoadManyReverse(tx, rec)
		if err != nil {
			return
		}
		recs = append(recs, newRecs...)
	}

	return
}

func (r *interfaceRelationship) Connect(tx RWTx, from, to Record) error {
	rel, err := r.getImplementingRelationship(tx, from.InterfaceID())
	if err != nil {
		return err
	}
	return rel.Connect(tx, from, to)
}

func (r *interfaceRelationship) Disconnect(tx RWTx, from, to Record) error {
	rel, err := r.getImplementingRelationship(tx, from.InterfaceID())
	if err != nil {
		return err
	}
	return rel.Disconnect(tx, from, to)
}
