package db

import (
	"fmt"
	"strings"
)

// // remove
// GetRelationships(Model) ([]Relationship, error)
// GetRelationship(ID) (Relationship, error)

// SaveModel(Model) error
// SaveRelationship(Relationship) error

type Schema struct {
	tx *holdTx
}

func (s *Schema) GetModelByID(mid ID) Model {
	mrec, err := s.tx.FindOne(ModelModel.ID(), EqID(ID(mid)))
	if err != nil {
		panic("GetModel failed")
	}
	return &model{mrec, s.tx}
}

func (s *Schema) GetModel(modelName string) (m Model, err error) {
	modelName = strings.ToLower(modelName)
	mrec, err := s.tx.h.FindOne(ModelModel.ID(), Eq("name", modelName))
	if err != nil {
		return m, fmt.Errorf("%w: %v", ErrInvalidModel, modelName)
	}
	return &model{mrec, s.tx}, nil
}

func (s *Schema) GetRelationship(id ID) (r Relationship, err error) {
	storeRel, err := s.tx.FindOne(RelationshipModel.ID(), EqID(id))
	if err != nil {
		return
	}
	return &rel{storeRel, s.tx}, nil
}

func (s *Schema) GetEnumValueByID(id ID) (ev EnumValue, err error) {
	storeEnumValue, err := s.tx.FindOne(EnumValueModel.ID(), EqID(id))
	if err != nil {
		return
	}
	return &enumValue{storeEnumValue}, nil
}

func (s *Schema) SaveRelationship(r Relationship) (err error) {
	rec, err := MarshalRecord(r, RelationshipModel)
	if err != nil {
		return
	}
	s.tx.Insert(rec)
	s.tx.Connect(rec.ID(), ID(r.Source().ID()), RelationshipSource)
	s.tx.Connect(rec.ID(), ID(r.Target().ID()), RelationshipTarget)
	return
}

func (s *Schema) SaveModel(m Model) (err error) {
	rec, err := MarshalRecord(m, ModelModel)
	if err != nil {
		return
	}

	s.tx.Insert(rec)
	attrs, _ := m.Attributes()
	for _, a := range attrs {
		var ar Record
		ar, err = MarshalRecord(a, AttributeModel)
		if err != nil {
			return
		}
		s.tx.Insert(ar)
		s.tx.Connect(rec.ID(), ar.ID(), ModelAttributes)
	}
	// done for side effect of gob registration
	RecordForModel(m)
	return
}
