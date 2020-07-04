package db

import (
	"fmt"
)

// Model

var ModelModel = ModelL{
	ID:         MakeID("872f8c55-9c12-43d1-b3f6-f7a02d937314"),
	Name:       "model",
	Attributes: []AttributeL{modelName},
}

var ModelAttributes = ConcreteRelationshipL{
	Name:   "attributes",
	ID:     MakeID("3271d6a5-0004-4752-81b8-b00142fd59bf"),
	Source: ModelModel,
	Target: ConcreteAttributeModel,
	Multi:  true,
}

var modelName = ConcreteAttributeL{
	Name:     "name",
	ID:       MakeID("d62d3c3a-0228-4131-98f5-2d49a2e3676a"),
	Datatype: String,
}

// Literal

type ModelL struct {
	ID         ID     `record:"id"`
	Name       string `record:"name"`
	Attributes []AttributeL
}

func (lit ModelL) GetID() ID {
	return lit.ID
}

func (lit ModelL) MarshalDB() (recs []Record, links []Link) {
	rec := MarshalRecord(lit, ModelModel)
	recs = append(recs, rec)
	for _, a := range lit.Attributes {
		ars, al := a.MarshalDB()
		recs = append(recs, ars...)
		links = append(links, al...)
		links = append(links, Link{rec.ID(), a.GetID(), ModelAttributes})
	}
	return
}

func (lit ModelL) AsModel() Model {
	return mlBox{lit}
}

// "Boxed" model

type mlBox struct {
	ModelL
}

func (m mlBox) ID() ID {
	return m.ModelL.ID
}

func (m mlBox) Lit() ModelL {
	return m.ModelL
}

func (m mlBox) Name() string {
	return m.ModelL.Name
}

func (m mlBox) Interfaces() ([]Interface, error) {
	panic("Not implemented")
}

func (m mlBox) Relationships() ([]Relationship, error) {
	panic("Not implemented")
}

func (m mlBox) RelationshipByName(name string) (Relationship, error) {
	panic("Not implemented")
}

func (m mlBox) Attributes() ([]Attribute, error) {
	var attrs []Attribute
	for _, a := range m.ModelL.Attributes {
		attrs = append(attrs, a.AsAttribute())
	}
	return attrs, nil
}

func (m mlBox) AttributeByName(name string) (a Attribute, err error) {
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

// Dynamic

type model struct {
	rec Record
	tx  Tx
}

func (m *model) ID() ID {
	return m.rec.ID()
}

func (m *model) Name() string {
	return modelName.AsAttribute().MustGet(m.rec).(string)
}

func (m *model) Relationships() (rels []Relationship, err error) {
	sourceRel, _ := m.tx.Schema().GetRelationshipByID(ConcreteRelationshipSource.ID)
	relRecs, err := m.tx.GetRelatedManyReverse(m.ID(), sourceRel)
	for _, rr := range relRecs {
		r := &concreteRelationship{rr, m.tx}
		rels = append(rels, r)
	}
	return
}

func (m *model) Attributes() (attrs []Attribute, err error) {
	attrRel, _ := m.tx.Schema().GetRelationshipByID(ConcreteRelationshipSource.ID)
	attrRecs, err := m.tx.GetRelatedMany(m.ID(), attrRel)
	for _, ar := range attrRecs {
		a := &concreteAttr{ar, m.tx}
		attrs = append(attrs, a)
	}
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
	a, ok := SystemAttrs[name]
	if !ok {
		err = fmt.Errorf("No attribute on model: %v %v", m.Name(), name)
	}
	return
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
