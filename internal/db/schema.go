package db

import (
	"fmt"
)

type Schema struct {
	tx *holdTx
	db *holdDB
}

func (s *Schema) GetInterfaceByID(id ID) (Interface, error) {
	irec, err := s.tx.h.FindOne(InterfaceInterface.ID(), EqID(id))
	if err != nil {
		return nil, err
	}
	return s.loadInterface(irec)
}

func (s *Schema) GetInterface(name string) (i Interface, err error) {
	ifaces := s.tx.Ref(InterfaceInterface.ID())
	irec, err := s.tx.Query(ifaces, Filter(ifaces, Eq("name", name))).OneRecord()
	if err != nil {
		return i, fmt.Errorf("%w: %v", ErrInvalidModel, name)
	}
	return s.loadInterface(irec)
}

func (s *Schema) GetModelByID(mid ID) (Model, error) {
	models := s.tx.Ref(ModelModel.ID())
	mrec, err := s.tx.Query(models, Filter(models, EqID(mid))).OneRecord()
	if err != nil {
		return nil, err
	}
	return s.LoadModel(mrec), nil
}

func (s *Schema) GetAttributeByID(attrID ID) (Attribute, error) {
	attrs := s.tx.Ref(ConcreteAttributeModel.ID())
	arec, err := s.tx.Query(attrs, Filter(attrs, EqID(attrID))).OneRecord()
	if err != nil {
		return nil, err
	}
	return s.LoadAttribute(arec), nil
}

func (s *Schema) GetDatatypeByID(datatypeID ID) (Datatype, error) {
	dts := s.tx.Ref(DatatypeInterface.ID())
	rec, err := s.tx.Query(dts, Filter(dts, EqID(datatypeID))).OneRecord()
	if err != nil {
		return nil, err
	}
	return s.loadDatatype(rec)
}

func (s *Schema) GetModel(modelName string) (m Model, err error) {
	models := s.tx.Ref(ModelModel.ID())
	mrec, err := s.tx.Query(models, Filter(models, Eq("name", modelName))).OneRecord()
	if err != nil {
		return m, fmt.Errorf("%w: %v", ErrInvalidModel, modelName)
	}
	return s.LoadModel(mrec), nil
}

func (s *Schema) GetRelationshipByID(id ID) (r Relationship, err error) {
	rels := s.tx.Ref(RelationshipInterface.ID())
	storeRel, err := s.tx.Query(rels, Filter(rels, EqID(id))).OneRecord()
	if err != nil {
		return nil, ErrNotFound
	}
	return s.loadRelationship(storeRel)
}

func (s *Schema) GetFunctionByID(id ID) (f Function, err error) {
	funcs := s.tx.Ref(FunctionInterface.ID())
	storeFunc, err := s.tx.Query(funcs, Filter(funcs, EqID(id))).OneRecord()
	if err != nil {
		return nil, ErrNotFound
	}
	return s.LoadFunction(storeFunc)
}

func (s *Schema) loadRelationship(rec Record) (Relationship, error) {
	interfaceID := rec.InterfaceID()
	rl, ok := s.db.rels[interfaceID]
	if !ok {
		return nil, ErrNotFound
	}
	return rl.Load(s.tx, rec), nil
}

func (s *Schema) LoadModel(rec Record) Model {
	return &model{rec, s.tx}
}

func (s *Schema) LoadAttribute(rec Record) Attribute {
	return &concreteAttr{rec, s.tx}
}

func (s *Schema) loadInterface(rec Record) (Interface, error) {
	il, ok := s.db.ifaces[rec.InterfaceID()]
	if !ok {
		return nil, ErrNotFound
	}
	return il.Load(s.tx, rec), nil
}

func (s *Schema) loadDatatype(rec Record) (Datatype, error) {
	dl, ok := s.db.datatypes[rec.InterfaceID()]
	if !ok {
		return nil, ErrNotFound
	}
	return dl.Load(s.tx, rec), nil
}

func (s *Schema) LoadFunction(rec Record) (Function, error) {
	fl, ok := s.db.runtimes[rec.InterfaceID()]
	if !ok {
		return nil, ErrNotFound
	}
	return fl.Load(s.tx, rec), nil
}

func (s *Schema) GetEnumValueByID(id ID) (ev EnumValue, err error) {
	evs := s.tx.Ref(EnumValueModel.ID())
	sev, err := s.tx.Query(evs, Filter(evs, EqID(id))).OneRecord()
	if err != nil {
		return
	}
	return &enumValue{sev}, nil
}
