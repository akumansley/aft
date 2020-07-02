package db

import (
	"fmt"
)

var ModelModel = ModelL{
	ID:         MakeID("872f8c55-9c12-43d1-b3f6-f7a02d937314"),
	Name:       "model",
	Attributes: []AttributeL{modelName},
}

var modelName = ConcreteAttributeL{
	Name:     "name",
	ID:       MakeID("d62d3c3a-0228-4131-98f5-2d49a2e3676a"),
	Datatype: String,
}

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
	sourceRel, _ := m.tx.Schema().GetRelationshipByID(RelationshipSource.ID)
	relRecs, err := m.tx.GetRelatedManyReverse(m.ID(), sourceRel)
	for _, rr := range relRecs {
		r := &rel{rr, m.tx}
		rels = append(rels, r)
	}
	return
}

func (m *model) Interfaces() (ifs []Interface, err error) {
	panic("not implemented")
}

func (m *model) Attributes() (attrs []Attribute, err error) {
	attrRel, _ := m.tx.Schema().GetRelationshipByID(RelationshipSource.ID)
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
		err = fmt.Errorf("No attribute on model: %v %v", m.Name, name)
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
	return nil, fmt.Errorf("No relationship on model: %v %v", m.Name, name)
}
