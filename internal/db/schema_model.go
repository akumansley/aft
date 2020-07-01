package db

import (
	"fmt"
)

type model struct {
	rec Record
	tx  Tx
}

func (m *model) ID() ID {
	return m.rec.ID()
}

func (m *model) Name() string {
	return m.rec.MustGet("name").(string)
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
		a := &attr{ar, m.tx}
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
