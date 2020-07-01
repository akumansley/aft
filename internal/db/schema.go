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
	mrec, err := s.tx.FindOne(ModelModel.ID, EqID(ID(mid)))
	if err != nil {
		panic("GetModel failed")
	}
	return &model{mrec, s.tx}
}

func (s *Schema) GetModel(modelName string) (m Model, err error) {
	modelName = strings.ToLower(modelName)
	mrec, err := s.tx.h.FindOne(ModelModel.ID, Eq("name", modelName))
	if err != nil {
		return m, fmt.Errorf("%w: %v", ErrInvalidModel, modelName)
	}
	return &model{mrec, s.tx}, nil
}

func (s *Schema) GetRelationshipByID(id ID) (r Relationship, err error) {
	storeRel, err := s.tx.FindOne(RelationshipModel.ID, EqID(id))
	if err != nil {
		return
	}
	return &rel{storeRel, s.tx}, nil
}

func (s *Schema) GetEnumValueByID(id ID) (ev EnumValue, err error) {
	storeEnumValue, err := s.tx.FindOne(EnumValueModel.ID, EqID(id))
	if err != nil {
		return
	}
	return &enumValue{storeEnumValue}, nil
}
