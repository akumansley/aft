package db

import (
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

var ModelRelationships = MakeConcreteRelationship(
	MakeID("420940ee-5745-429c-bc10-3e43ec8b9a63"),
	"relationships",
	true,
	RelationshipInterface,
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

func (l ModelInterfaceLoader) Load(tx Tx, rec Record) Interface {
	return &model{rec, tx}
}

// Literal

func MakeModel(id ID, name string, attrs []AttributeL, rels []RelationshipL, implements []ConcreteInterfaceL) ModelL {
	return ModelL{
		id, name, attrs, rels, implements,
	}
}

type ModelL struct {
	ID_            ID     `record:"id"`
	Name_          string `record:"name"`
	Attributes_    []AttributeL
	Relationships_ []RelationshipL
	Implements_    []ConcreteInterfaceL
}

func (lit ModelL) MarshalDB() (recs []Record, links []Link) {
	rec := MarshalRecord(lit, ModelModel)
	rec.Set("system", true)
	recs = append(recs, rec)
	for _, a := range lit.Attributes_ {
		ars, al := a.MarshalDB()
		recs = append(recs, ars...)
		links = append(links, al...)
		links = append(links, Link{rec.ID(), a.ID(), ModelAttributes})
	}

	for _, r := range lit.Relationships_ {
		rrecs, rlinks := r.MarshalDB()
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

func (lit ModelL) Name() string {
	return lit.Name_
}

func (lit ModelL) Interfaces() ([]Interface, error) {
	panic("Not implemented")
}

func (lit ModelL) Relationships() ([]Relationship, error) {
	panic("Not implemented")
}

func (lit ModelL) RelationshipByName(name string) (Relationship, error) {
	panic("Not implemented")
}

func (lit ModelL) Attributes() ([]Attribute, error) {
	var attrs []Attribute
	for _, a := range lit.Attributes_ {
		attrs = append(attrs, a)
	}
	attrs = append(attrs, MakeConcreteAttribute(lit.ID_, "id", UUID))
	return attrs, nil
}

func (lit ModelL) Implements() ([]Interface, error) {
	var ifaces []Interface
	for _, i := range lit.Implements_ {
		ifaces = append(ifaces, i)
	}
	return ifaces, nil
}

func (lit ModelL) AttributeByName(name string) (a Attribute, err error) {
	attrs, err := lit.Attributes()
	if err != nil {
		return
	}
	for _, attr := range attrs {
		if attr.Name() == name {
			return attr, nil
		}
	}
	return nil, fmt.Errorf("No attribute on model: %v %v", lit.Name(), name)
}

// Dynamic

type model struct {
	rec Record
	tx  Tx
}

func (m *model) ID() ID {
	return m.rec.ID()
}

func (m *model) Name() string {
	return modelName.MustGet(m.rec).(string)
}

func (m *model) Relationships() (rels []Relationship, err error) {
	relRecs, err := m.tx.getRelatedMany(m.ID(), ModelRelationships.ID())
	if err != nil {
		return
	}
	for _, rr := range relRecs {
		r, err := m.tx.Schema().loadRelationship(rr)
		if err != nil {
			return nil, err
		}
		rels = append(rels, r)
	}
	return
}

func (m *model) Attributes() (attrs []Attribute, err error) {
	attrRecs, err := m.tx.getRelatedMany(m.ID(), ModelAttributes.ID())
	if err != nil {
		return
	}

	for _, ar := range attrRecs {
		a := &concreteAttr{ar, m.tx}
		attrs = append(attrs, a)
	}
	// refactor..
	attrs = append(attrs, MakeConcreteAttribute(m.ID(), "id", UUID))
	attrs = append(attrs, MakeConcreteAttribute(m.ID(), "type", Type))
	return
}

// TODO rewrite as a findone
func (m *model) AttributeByName(name string) (a Attribute, err error) {
	attrs, err := m.Attributes()
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
func (m *model) RelationshipByName(name string) (rel Relationship, err error) {
	rels, err := m.Relationships()
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

func (m *model) Implements() (ifs []Interface, err error) {
	ifRecs, err := m.tx.getRelatedMany(m.ID(), ModelImplements.ID())
	if err != nil {
		return
	}

	for _, ir := range ifRecs {
		i := &iface{ir, m.tx}
		ifs = append(ifs, i)
	}
	return
}
