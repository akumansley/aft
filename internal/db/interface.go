package db

import (
	"fmt"
)

// Interface

var InterfaceInterface = MakeInterface(
	MakeID("7858a890-bf53-49fe-8ef3-145b6a43bc4c"),
	"interface",
	[]AttributeL{interfaceInterfaceName},
	[]InterfaceRelationshipL{
		AbstractInterfaceAttributes,
		AbstractInterfaceRelationships,
		AbstractInterfaceModule,
	},
)

var interfaceInterfaceName = MakeConcreteAttribute(
	MakeID("54be2897-8f4f-4906-aba8-91bdc430275e"),
	"name",
	String,
)

var AbstractInterfaceRelationships = MakeInterfaceRelationship(
	MakeID("f94881af-de80-4bb9-8199-3d33ff248280"),
	"relationships",
	true,
	RelationshipInterface,
)

var AbstractInterfaceAttributes = MakeInterfaceRelationship(
	MakeID("a910aa8d-b8fc-47d7-ab44-0d5f5607dad9"),
	"attributes",
	true,
	ConcreteAttributeModel,
)

var AbstractInterfaceModule = MakeInterfaceRelationship(
	MakeID("5512ceee-9e1a-44e2-871b-19fd6f4d3e50"),
	"module",
	false,
	ModuleModel,
)

var abstractInterfaceName = MakeConcreteAttribute(
	MakeID("b70c2d2a-a9ec-4e70-b6f8-7c9d3beb2419"),
	"name",
	String,
)

// Model

var InterfaceModel = MakeModel(
	MakeID("a9bab408-fb98-463c-a6e3-4613adb8dca4"),
	"concreteInterface",
	[]AttributeL{interfaceName},
	[]RelationshipL{
		InterfaceAttributes,
		InterfaceRelationships,
		InterfaceModule,
	},
	[]ConcreteInterfaceL{InterfaceInterface},
)

var InterfaceRelationships = MakeConcreteRelationship(
	MakeID("485cfc71-3941-4458-979d-185f10a225b2"),
	"relationships",
	true,
	InterfaceRelationshipModel,
)

var InterfaceAttributes = MakeConcreteRelationship(
	MakeID("cf534a84-852a-40d5-b5cf-8457db120e58"),
	"attributes",
	true,
	ConcreteAttributeModel,
)

var InterfaceModule = MakeConcreteRelationship(
	MakeID("33b037c5-800d-41f3-baf1-2cb4ef431153"),
	"module",
	false,
	ModuleModel,
)

var interfaceName = MakeConcreteAttribute(
	MakeID("f3064600-5a9e-45ce-b832-0e25d9c18434"),
	"name",
	String,
)

// Loader
type InterfaceInterfaceLoader struct{}

func (l InterfaceInterfaceLoader) ProvideModel() ModelL {
	return InterfaceModel
}

func (l InterfaceInterfaceLoader) Load(rec Record) Interface {
	return &iface{rec}
}

// Literal

func MakeInterface(id ID, name string, attrs []AttributeL, rels []InterfaceRelationshipL) ConcreteInterfaceL {
	return ConcreteInterfaceL{
		id, name, attrs, rels,
	}
}

type ConcreteInterfaceL struct {
	ID_            ID     `record:"id"`
	Name_          string `record:"name"`
	Attributes_    []AttributeL
	Relationships_ []InterfaceRelationshipL
}

func (lit ConcreteInterfaceL) MarshalDB(b *Builder) (recs []Record, links []Link) {
	rec := MarshalRecord(b, lit)
	recs = append(recs, rec)
	for _, a := range lit.Attributes_ {
		ars, al := a.MarshalDB(b)
		recs = append(recs, ars...)
		links = append(links, al...)
		links = append(links, Link{lit, a, InterfaceAttributes})
	}

	for _, r := range lit.Relationships_ {
		rrecs, rlinks := r.MarshalDB(b)
		recs = append(recs, rrecs...)
		links = append(links, rlinks...)
		links = append(links, Link{lit, r, InterfaceRelationships})
	}
	return
}

func (lit ConcreteInterfaceL) ID() ID {
	return lit.ID_
}

func (lit ConcreteInterfaceL) InterfaceID() ID {
	return InterfaceModel.ID()
}

func (lit ConcreteInterfaceL) InterfaceName() string {
	return InterfaceModel.Name_
}

func (lit ConcreteInterfaceL) Load(tx Tx) Interface {
	iface, err := tx.Schema().GetInterfaceByID(lit.ID())
	if err != nil {
		panic(err)
	}
	return iface
}

// Dynamic

type iface struct {
	rec Record
}

func (m *iface) ID() ID {
	return m.rec.ID()
}

func (m *iface) Name() string {
	return m.rec.MustGet("name").(string)
}

func (m *iface) Relationships(tx Tx) (rels []Relationship, err error) {
	relRecs, err := tx.getRelatedMany(m.ID(), InterfaceRelationships.ID())
	if err != nil {
		return
	}

	// TODO is this correct?
	for _, rr := range relRecs {
		var r Relationship
		r, err = tx.Schema().loadRelationship(rr)
		if err != nil {
			return
		}
		rels = append(rels, r)
	}
	return
}

func (m *iface) Attributes(tx Tx) (attrs []Attribute, err error) {
	attrRecs, err := tx.getRelatedMany(m.ID(), InterfaceAttributes.ID())
	if err != nil {
		return
	}

	for _, ar := range attrRecs {
		a := &concreteAttr{ar}
		attrs = append(attrs, a)
	}
	id, err := tx.Schema().GetAttributeByID(GlobalIDAttribute.ID())
	if err != nil {
		return
	}
	type_, _ := tx.Schema().GetAttributeByID(GlobalTypeAttribute.ID())
	if err != nil {
		return
	}
	attrs = append(attrs, type_, id)
	return
}

// TODO rewrite as a findone
func (m *iface) AttributeByName(tx Tx, name string) (a Attribute, err error) {
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
func (m *iface) RelationshipByName(tx Tx, name string) (rel Relationship, err error) {
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
