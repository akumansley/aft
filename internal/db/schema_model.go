package db

import (
	"fmt"
)

type model struct {
	rec Record
	tx  Tx
}

func (m *model) ID() ModelID {
	return ModelID(m.rec.ID())
}

func (m *model) Name() string {
	return m.rec.MustGet("name").(string)
}

func (m *model) Relationships() (rels []Relationship, err error) {
	relRecs, err := m.tx.GetRelatedManyReverse(ID(m.ID()), RelationshipSource.ID)
	for _, rr := range relRecs {
		r := &rel{rr, m.tx}
		rels = append(rels, r)
	}
	return
}

func (m *model) Attributes() (attrs []Attribute, err error) {
	attrRecs, err := m.tx.GetRelatedMany(ID(m.ID()), ModelAttributes.ID())
	for _, ar := range attrRecs {
		a := &attr{ar, m.tx}
		attrs = append(attrs, a)
	}
	return
}

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
