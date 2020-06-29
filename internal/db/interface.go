package db

import (
	"fmt"
)

// Interface

var InterfaceInterface = MakeInterface(
	MakeID("7858a890-bf53-49fe-8ef3-145b6a43bc4c"),
	"interface",
	[]AttributeL{interfaceName},
	[]RelationshipL{
		InterfaceAttributes,
		InterfaceRelationships,
	},
)

var AbstractInterfaceRelationships = MakeConcreteRelationship(
	MakeID("485cfc71-3941-4458-979d-185f10a225b2"),
	"relationships",
	true,
	ConcreteRelationshipModel,
)

var AbstractInterfaceAttributes = MakeConcreteRelationship(
	MakeID("a910aa8d-b8fc-47d7-ab44-0d5f5607dad9"),
	"attributes",
	true,
	ConcreteAttributeModel,
)

var abstractInterfaceName = MakeConcreteAttribute(
	MakeID("b70c2d2a-a9ec-4e70-b6f8-7c9d3beb2419"),
	"name",
	String,
)

// Model

var InterfaceModel = MakeModel(
	MakeID("a9bab408-fb98-463c-a6e3-4613adb8dca4"),
	"concreteinterface",
	[]AttributeL{interfaceName},
	[]RelationshipL{
		InterfaceAttributes,
		InterfaceRelationships,
	},
	[]ConcreteInterfaceL{InterfaceInterface},
)

var InterfaceRelationships = MakeConcreteRelationship(
	MakeID("485cfc71-3941-4458-979d-185f10a225b2"),
	"relationships",
	true,
	ConcreteRelationshipModel,
)

var InterfaceAttributes = MakeConcreteRelationship(
	MakeID("cf534a84-852a-40d5-b5cf-8457db120e58"),
	"attributes",
	true,
	ConcreteAttributeModel,
)

var interfaceName = MakeConcreteAttribute(
	MakeID("f3064600-5a9e-45ce-b832-0e25d9c18434"),
	"name",
	String,
)

// Literal

func MakeInterface(id ID, name string, attrs []AttributeL, rels []RelationshipL) ConcreteInterfaceL {
	return ConcreteInterfaceL{
		id, name, attrs, rels,
	}
}

type ConcreteInterfaceL struct {
	ID_            ID     `record:"id"`
	Name_          string `record:"name"`
	Attributes_    []AttributeL
	Relationships_ []RelationshipL
}

func (lit ConcreteInterfaceL) MarshalDB() (recs []Record, links []Link) {
	rec := MarshalRecord(lit, InterfaceModel)
	recs = append(recs, rec)
	for _, a := range lit.Attributes_ {
		ars, al := a.MarshalDB()
		recs = append(recs, ars...)
		links = append(links, al...)
		links = append(links, Link{rec.ID(), a.ID(), InterfaceAttributes})
	}

	for _, r := range lit.Relationships_ {
		rrecs, rlinks := r.MarshalDB()
		recs = append(recs, rrecs...)
		links = append(links, rlinks...)
		switch r.(type) {
		case ConcreteRelationshipL:
			links = append(links, Link{rec.ID(), r.ID(), InterfaceRelationships})
		}
	}
	return
}

func (lit ConcreteInterfaceL) ID() ID {
	return lit.ID_
}

func (lit ConcreteInterfaceL) Name() string {
	return lit.Name_
}

func (lit ConcreteInterfaceL) Relationships() ([]Relationship, error) {
	panic("Not implemented")
}

func (lit ConcreteInterfaceL) RelationshipByName(name string) (Relationship, error) {
	panic("Not implemented")
}

func (lit ConcreteInterfaceL) Attributes() ([]Attribute, error) {
	var attrs []Attribute
	for _, a := range lit.Attributes_ {
		attrs = append(attrs, a)
	}
	return attrs, nil
}

func (lit ConcreteInterfaceL) AttributeByName(name string) (a Attribute, err error) {
	attrs, err := lit.Attributes()
	if err != nil {
		return
	}
	for _, attr := range attrs {
		if attr.Name() == name {
			return attr, nil
		}
	}
	a, ok := SystemAttrs[name]
	if !ok {
		err = fmt.Errorf("No attribute on interface: %v %v", lit.Name(), name)
	}
	return
}

// Dynamic

type iface struct {
	rec Record
	tx  Tx
}

func (m *iface) ID() ID {
	return m.rec.ID()
}

func (m *iface) Name() string {
	return interfaceName.MustGet(m.rec).(string)
}

func (m *iface) Relationships() (rels []Relationship, err error) {
	relRecs, err := m.tx.getRelatedMany(m.ID(), InterfaceRelationships.ID())
	if err != nil {
		return
	}
	for _, rr := range relRecs {
		r := &concreteRelationship{rr, m.tx}
		rels = append(rels, r)
	}
	return
}

func (m *iface) Attributes() (attrs []Attribute, err error) {
	attrRecs, err := m.tx.getRelatedMany(m.ID(), InterfaceAttributes.ID())
	if err != nil {
		return
	}

	for _, ar := range attrRecs {
		a := &concreteAttr{ar, m.tx}
		attrs = append(attrs, a)
	}
	return
}

// TODO rewrite as a findone
func (m *iface) AttributeByName(name string) (a Attribute, err error) {
	attrs, err := m.Attributes()
	if err != nil {
		return
	}
	for _, attr := range attrs {
		if attr.Name() == name {
			return attr, nil
		}
	}
	a, ok := SystemAttrs[name]
	if !ok {
		err = fmt.Errorf("No attribute on model: %v %v", m.Name(), name)
	}
	return
}

// TODO rewrite as a findone
func (m *iface) RelationshipByName(name string) (rel Relationship, err error) {
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
