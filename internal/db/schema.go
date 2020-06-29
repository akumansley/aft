package db

import (
	"fmt"
)

// // remove
// GetModel(string) (Model, error)
// GetRelationships(Model) ([]Relationship, error)
// GetRelationship(ID) (Relationship, error)
// GetModelByID(ModelID) (Model, error)
// SaveModel(Model) error
// SaveRelationship(Relationship) error

type Schema struct {
	tx Tx
}

type Model interface {
	ID() ModelID
	Name() string
	Attributes() []Attribute
	Relationships() []Relationship
}

func (s *Schema) GetModel(mid ModelID) Model {
	mrec, err := s.tx.FindOne(ModelModel.ID, ID(mid))
	if err != nil {
		panic("GetModel failed")
	}
	return &model{mrec, tx}
}

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
	attrRecs, err := m.tx.GetRelatedMany(ID(m.ID()), ModelAttributes.ID)
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

type Attribute interface {
	ID() ID
	Name() string
	// Datatype() Datatype
}

type attr struct {
	rec Record
	tx  Tx
}

func (a *attr) ID() ID {
	return a.rec.ID()
}

func (a *attr) Name() string {
	return a.rec.MustGet("name").(string)
}

type Relationship interface {
	ID() ID
	Name() string
	Multi() bool
	Source() Model
	Target() Model
}

type rel struct {
	rec Record
	tx  Tx
}

func (r *rel) ID() ID {
	return r.rec.ID()
}

func (r *rel) Name() string {
	return r.rec.MustGet("name").(string)
}

func (r *rel) Multi() bool {
	return r.rec.MustGet("multi").(bool)
}

func (r *rel) Source() Model {
	mRec, err := r.tx.GetRelatedOne(r.ID(), RelationshipSource)
	if err != nil {
		panic("source failed")
	}
	return &model{mRec, tx}
}

func (r *rel) Source() Model {
	mRec, err := r.tx.GetRelatedOne(r.ID(), RelationshipTarget)
	if err != nil {
		panic("source failed")
	}
	return &model{mRec, tx}
}
