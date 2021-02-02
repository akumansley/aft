package db

import (
	"encoding/gob"
	"fmt"
)

// Model

var ModelModel = MakeModel(
	MakeID("872f8c55-9c12-43d1-b3f6-f7a02d937314"),
	"model",
	[]AttributeL{modelName,
		modelSystem},
	[]RelationshipL{
		ModelAttributes,
		ModelRelationships,
		ModelImplements,
	},
	[]ConcreteInterfaceL{InterfaceInterface},
)

var GlobalIDAttribute = MakeConcreteAttribute(
	MakeID("ea8e3b18-7723-4005-b357-4384ae87fdaa"),
	"id",
	UUID,
)

var GlobalTypeAttribute = MakeConcreteAttribute(
	MakeID("fbf4a0c4-2766-4708-aa56-2b1620e0e15a"),
	"type",
	Type,
)

var ModelRelationships = MakeConcreteRelationship(
	MakeID("420940ee-5745-429c-bc10-3e43ec8b9a63"),
	"relationships",
	true,
	RelationshipInterface,
)

var ModelTargeted = MakeReverseRelationship(
	MakeID("f0a8f141-3157-4b43-b484-1e784f7c69da"),
	"targeted",
	ConcreteRelationshipTarget,
)

var ModelAttributes = MakeConcreteRelationship(
	MakeID("3271d6a5-0004-4752-81b8-b00142fd59bf"),
	"attributes",
	true,
	ConcreteAttributeModel,
)

var ModelImplements = MakeConcreteRelationship(
	MakeID("0b1e45e4-7a68-435c-9e53-b8ba0cff5f5d"),
	"implements",
	true,
	InterfaceModel,
)

var modelName = MakeConcreteAttribute(
	MakeID("d62d3c3a-0228-4131-98f5-2d49a2e3676a"),
	"name",
	String,
)

var modelSystem = MakeConcreteAttribute(
	MakeID("2c7a33f8-0baf-4e02-b584-7681095d6c2e"),
	"system",
	Bool,
)

// Loader
type ModelInterfaceLoader struct{}

func (l ModelInterfaceLoader) ProvideModel() ModelL {
	return ModelModel
}

func (l ModelInterfaceLoader) Load(rec Record) Interface {
	return &model{rec}
}

// Literal

func MakeModel(id ID, name string, attrs []AttributeL, rels []RelationshipL, implements []ConcreteInterfaceL) ModelL {
	return ModelL{
		id, name, true, attrs, rels, implements,
	}
}

type ModelL struct {
	ID_            ID     `record:"id"`
	Name_          string `record:"name"`
	System         bool   `record:"system"`
	Attributes_    []AttributeL
	Relationships_ []RelationshipL
	Implements_    []ConcreteInterfaceL
}

func (lit ModelL) MarshalDB(b *Builder) (recs []Record, links []Link) {
	rec := MarshalRecord(b, lit)
	recs = append(recs, rec)
	for _, a := range lit.Attributes_ {
		ars, al := a.MarshalDB(b)
		recs = append(recs, ars...)
		links = append(links, al...)
		links = append(links, Link{rec.ID(), a.ID(), ModelAttributes})
	}

	for _, r := range lit.Relationships_ {
		rrecs, rlinks := r.MarshalDB(b)
		recs = append(recs, rrecs...)
		links = append(links, rlinks...)
		links = append(links, Link{rec.ID(), r.ID(), ModelRelationships})
	}

	for _, i := range lit.Implements_ {
		links = append(links, Link{rec.ID(), i.ID(), ModelImplements})
	}
	return
}

func (lit ModelL) ID() ID {
	return lit.ID_
}

func (lit ModelL) InterfaceID() ID {
	return ModelModel.ID()
}

func (lit ModelL) InterfaceName() string {
	return ModelModel.Name_
}

func (lit ModelL) Load(tx Tx) Interface {
	iface, err := tx.Schema().GetInterfaceByID(lit.ID())
	if err != nil {
		panic(err)
	}
	return iface
}

// Dynamic

func init() {
	gob.Register(&model{})
}

type model struct {
	Rec Record
}

func (m *model) ID() ID {
	return m.Rec.ID()
}

func (m *model) Name() string {
	return m.Rec.MustGet("name").(string)
}

func (m *model) Relationships(tx Tx) (rels []Relationship, err error) {
	relRecs, err := tx.getRelatedMany(m.ID(), ModelRelationships.ID())
	if err != nil {
		return
	}
	for _, rr := range relRecs {
		r, err := tx.Schema().loadRelationship(rr)
		if err != nil {
			return nil, err
		}
		rels = append(rels, r)
	}
	return
}

func (m *model) Targeted(tx Tx) (rels []Relationship, err error) {
	relRecs, err := tx.getRelatedMany(m.ID(), ModelTargeted.ID())
	if err != nil {
		return
	}
	for _, rr := range relRecs {
		r, err := tx.Schema().loadRelationship(rr)
		if err != nil {
			return nil, err
		}
		rels = append(rels, r)
	}
	return
}

func (m *model) Attributes(tx Tx) (attrs []Attribute, err error) {
	attrRecs, err := tx.getRelatedMany(m.ID(), ModelAttributes.ID())
	if err != nil {
		return
	}

	for _, ar := range attrRecs {
		a := &concreteAttr{ar}
		attrs = append(attrs, a)
	}
	// refactor..
	id, err := tx.Schema().GetAttributeByID(GlobalIDAttribute.ID())
	if err != nil {
		return
	}
	type_, err := tx.Schema().GetAttributeByID(GlobalTypeAttribute.ID())
	if err != nil {
		return
	}
	attrs = append(attrs, type_, id)
	return
}

// TODO rewrite as a findone
func (m *model) AttributeByName(tx Tx, name string) (a Attribute, err error) {
	attrs, err := m.Attributes(tx)
	if err != nil {
		return
	}
	for _, attr := range attrs {
		if attr.Name() == name {
			return attr, nil
		}
	}
	return nil, fmt.Errorf("No attribute on model: %v %v", m.Name(), name)
}

// TODO rewrite as a findone
func (m *model) RelationshipByName(tx Tx, name string) (rel Relationship, err error) {
	rels, err := m.Relationships(tx)
	if err != nil {
		return
	}
	for _, rel := range rels {
		if rel.Name() == name {
			return rel, nil
		}
	}
	return nil, fmt.Errorf("No relationship on model: %v %v", m.Name(), name)
}

func (m *model) Implements(tx Tx) (ifs []Interface, err error) {
	ifRecs, err := tx.getRelatedMany(m.ID(), ModelImplements.ID())
	if err != nil {
		return
	}

	for _, ir := range ifRecs {
		i := &iface{ir}
		ifs = append(ifs, i)
	}
	return
}
