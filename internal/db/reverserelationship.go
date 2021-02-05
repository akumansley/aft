package db

import (
	"encoding/gob"
	"fmt"
)

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
	RelationshipInterface,
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

func MakeReverseRelationship(id ID, name string, referencing RelationshipL) ReverseRelationshipL {
	return ReverseRelationshipL{
		id, name, referencing, nil,
	}
}

func MakeReverseRelationshipWithSource(id ID, name string, referencing RelationshipL, source InterfaceL) ReverseRelationshipL {
	return ReverseRelationshipL{
		id, name, referencing, source,
	}
}

type ReverseRelationshipL struct {
	ID_          ID     `record:"id"`
	Name_        string `record:"name"`
	Referencing_ RelationshipL
	Source       InterfaceL
}

func (lit ReverseRelationshipL) MarshalDB(b *Builder) (recs []Record, links []Link) {
	rec := MarshalRecord(b, lit)
	recs = append(recs, rec)
	links = []Link{
		Link{lit, lit.Referencing_, ReverseRelationshipReferencing},
	}
	if lit.Source != nil {
		links = append(links, Link{lit.Source, lit, ModelRelationships})
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

func init() {
	gob.Register(&reverseRelationship{})
}

type reverseRelationship struct {
	Rec Record
}

func (r *reverseRelationship) ID() ID {
	return r.Rec.ID()
}

func (r *reverseRelationship) Name() string {
	return r.Rec.MustGet("name").(string)
}

func (r *reverseRelationship) Multi() bool {
	// until we have constraints, we assume it's true
	return true
}

func (r *reverseRelationship) getReferencing(tx Tx) Relationship {
	referenced, err := tx.getRelatedOne(r.ID(), ReverseRelationshipReferencing.ID())
	if err != nil {
		err := fmt.Errorf("rev %v referencing failed \n", r.ID())
		panic(err)
	}
	referencedRel, err := tx.Schema().loadRelationship(referenced)
	if err != nil {
		err := fmt.Errorf("rev %v referencing failed \n", r.ID())
		panic(err)
	}
	if rev, ok := referencedRel.(*reverseRelationship); ok {
		err := fmt.Errorf("reverse rel %v referencing reverse rel %v is invalid", r, rev)
		panic(err)
	}
	return referencedRel

}

func (r *reverseRelationship) Source(tx Tx) Interface {
	referencedRel := r.getReferencing(tx)
	return referencedRel.Target(tx)
}

func (r *reverseRelationship) Target(tx Tx) Interface {
	referencedRel := r.getReferencing(tx)
	return referencedRel.Source(tx)
}

func (r *reverseRelationship) LoadOne(tx Tx, rec Record) (Record, error) {
	panic("Not implemented")
}

func (r *reverseRelationship) LoadMany(tx Tx, rec Record) ([]Record, error) {
	referenced := r.getReferencing(tx)
	if !r.Multi() {
		panic("LoadMany on non-multi record")
	}
	return referenced.LoadManyReverse(tx, rec)
}

func (r *reverseRelationship) LoadManyReverse(tx Tx, rec Record) ([]Record, error) {
	referenced := r.getReferencing(tx)
	if !r.Multi() {
		panic("LoadMany on non-multi record")
	}
	return referenced.LoadMany(tx, rec)
}

func (r *reverseRelationship) Connect(tx RWTx, c, p Record) error {
	referenced := r.getReferencing(tx)
	return referenced.Connect(tx, p, c)
}

func (r *reverseRelationship) Disconnect(tx RWTx, c, p Record) error {
	referenced := r.getReferencing(tx)
	return referenced.Disconnect(tx, p, c)
}
