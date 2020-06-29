package db

import (
	"fmt"
	"strings"
)

type Schema struct {
	tx *holdTx
}

func (s *Schema) GetModelByID(mid ID) (Model, error) {
	mrec, err := s.tx.FindOne(ModelModel.ID(), EqID(mid))
	if err != nil {
		return nil, err
	}
	return &model{mrec, s.tx}, nil
}

func (s *Schema) GetModel(modelName string) (m Model, err error) {
	modelName = strings.ToLower(modelName)
	mrec, err := s.tx.h.FindOne(ModelModel.ID(), Eq("name", modelName))
	if err != nil {
		return m, fmt.Errorf("%w: %v", ErrInvalidModel, modelName)
	}
	return &model{mrec, s.tx}, nil
}

func (s *Schema) GetRelationshipByID(id ID) (r Relationship, err error) {
	for mid, rl := range s.tx.db.rels {
		storeRel, err := s.tx.FindOne(mid, EqID(id))
		if err != nil {
			continue
		}
		r = rl.Load(s.tx, storeRel)
	}
	return r, err
}

func (s *Schema) loadRelationship(rec Record) (Relationship, error) {
	iface := rec.Interface()
	rl, ok := s.tx.db.rels[iface.ID()]
	if !ok {
		return nil, ErrNotFound
	}
	return rl.Load(s.tx, rec), nil
}

func (s *Schema) LoadFunction(rec Record) (Function, error) {
	iface := rec.Interface()
	fl, ok := s.tx.db.runtimes[iface.ID()]
	if !ok {
		return nil, ErrNotFound
	}
	return fl.Load(s.tx, rec), nil
}

func (s *Schema) GetEnumValueByID(id ID) (ev EnumValue, err error) {
	storeEnumValue, err := s.tx.FindOne(EnumValueModel.ID(), EqID(id))
	if err != nil {
		return
	}
	return &enumValue{storeEnumValue}, nil
}
