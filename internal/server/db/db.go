package db

import (
	"awans.org/aft/er"
	"awans.org/aft/er/q"
	"awans.org/aft/internal/model"
	"errors"
	"fmt"
	"strings"
	"sync"
)

var (
	ErrData         = errors.New("data-error")
	ErrInvalidModel = fmt.Errorf("%w: invalid model", ErrData)
)

func New() DB {
	appDB := holdDB{h: er.New()}
	appDB.AddMetaModel()
	return &appDB
}

func (db *holdDB) AddMetaModel() {
	tx := db.NewRWTx()
	tx.SaveModel(ModelModel)
	tx.SaveModel(AttributeModel)
	tx.SaveModel(RelationshipModel)
	tx.Commit()
}

// DB is a value
type DB interface {
	NewTx() Tx
	NewRWTx() RWTx
}

type Tx interface {
	GetModel(string) (model.Model, error)
	MakeRecord(string) model.Record
	Resolve(*model.IncludeResult, Inclusion)
	FindOne(string, UniqueQuery) (model.Record, error)
	FindMany(string, Query) []model.Record
}

type RWTx interface {
	GetModel(string) (model.Model, error)
	MakeRecord(string) model.Record
	Resolve(*model.IncludeResult, Inclusion)
	FindOne(string, UniqueQuery) (model.Record, error)
	FindMany(string, Query) []model.Record

	SaveModel(model.Model)
	Insert(model.Record)
	Connect(from, to model.Record, fromRel model.Relationship)
	Commit()
}

type holdDB struct {
	sync.RWMutex
	h *er.Hold
}

type holdTx struct {
	h  *er.Hold
	db *holdDB
	rw bool
}

func (tx *holdTx) ensureWrite() {
	if !tx.rw {
		panic("Tried to write in a read only tx")
	}
}

func (db *holdDB) NewTx() Tx {
	db.RLock()
	tx := holdTx{h: db.h, rw: false}
	db.RUnlock()
	return &tx
}

func (db *holdDB) NewRWTx() RWTx {
	db.RLock()
	tx := holdTx{h: db.h, db: db, rw: true}
	db.RUnlock()
	return &tx
}

func (tx *holdTx) FindOne(modelName string, uq UniqueQuery) (rec model.Record, err error) {
	rec, err = tx.h.FindOne(modelName, q.Eq(uq.Key, uq.Val))
	return
}

func (tx *holdTx) Insert(rec model.Record) {
	tx.ensureWrite()
	tx.h = tx.h.Insert(rec)
}

// TODO hack -- remove this and rewrite with Relationship containing the name
func getBackref(tx Tx, rel model.Relationship) model.Relationship {
	m, _ := tx.GetModel(rel.TargetModel)
	return m.Relationships[rel.TargetRel]
}

func (tx *holdTx) Connect(from, to model.Record, fromRel model.Relationship) {
	tx.ensureWrite()
	toRel := getBackref(tx, fromRel)
	if fromRel.RelType == model.BelongsTo && (toRel.RelType == model.HasOne || toRel.RelType == model.HasMany) {
		// FK from
		from.SetFK(toRel.TargetRel, to.Id())
	} else if toRel.RelType == model.BelongsTo && (fromRel.RelType == model.HasOne || fromRel.RelType == model.HasMany) {
		// FK to
		to.SetFK(fromRel.TargetRel, from.Id())
	} else if toRel.RelType == model.HasManyAndBelongsToMany && fromRel.RelType == model.HasManyAndBelongsToMany {
		// Join table
		panic("Many to many relationships not implemented yet")
	} else {
		fmt.Printf("fromRel %v toRel %v\n", fromRel, toRel)
		panic("Trying to connect invalid relationship")
	}
	h1 := tx.h.Insert(from)
	h2 := h1.Insert(to)
	tx.h = h2
}

func (tx *holdTx) GetModel(modelName string) (m model.Model, err error) {
	modelName = strings.ToLower(modelName)
	ifc, err := tx.h.FindOne("model", q.Eq("name", modelName))
	if err != nil {
		return m, fmt.Errorf("%w: %v", ErrInvalidModel, modelName)
	}
	storeModel := ifc.(model.Record)

	m = model.Model{
		Type: storeModel.Type(),
		Id:   storeModel.Id(),
		Name: storeModel.Get("name").(string),
	}

	attrs := make(map[string]model.Attribute)

	// make ModelId a dynamic key
	ami := tx.h.IterMatches("attribute", q.EqFK("model", m.Id))
	for storeAttrIf, ok := ami.Next(); ok; storeAttrIf, ok = ami.Next() {
		storeAttr := storeAttrIf.(model.Record)
		attr := model.Attribute{
			AttrType: model.AttrType(storeAttr.Get("attrType").(int64)),
			Type:     storeAttr.Type(),
			Id:       storeAttr.Id(),
		}
		name := storeAttr.Get("name").(string)
		attrs[name] = attr
	}
	m.Attributes = attrs

	rels := make(map[string]model.Relationship)

	// make ModelId a dynamic key
	rmi := tx.h.IterMatches("relationship", q.EqFK("model", m.Id))
	for storeRelIf, ok := rmi.Next(); ok; storeRelIf, ok = rmi.Next() {
		storeRel := storeRelIf.(model.Record)
		rel := model.Relationship{
			Type:        storeRel.Type(),
			Id:          storeRel.Id(),
			RelType:     model.RelType(storeRel.Get("relType").(int64)),
			TargetModel: storeRel.Get("targetModel").(string),
			TargetRel:   storeRel.Get("targetRel").(string),
		}
		name := storeRel.Get("name").(string)
		rels[name] = rel
	}
	m.Relationships = rels

	return m, nil
}

// Manual serialization required for bootstrapping
func (tx *holdTx) SaveModel(m model.Model) {
	tx.ensureWrite()
	storeModel := model.RecordForModel(ModelModel)
	storeModel.Set("name", m.Name)
	storeModel.Set("id", m.Id)
	tx.h = tx.h.Insert(storeModel)

	for aKey, attr := range m.Attributes {
		storeAttr := model.RecordForModel(AttributeModel)
		storeAttr.Set("name", aKey)
		storeAttr.Set("attrType", int64(attr.AttrType))
		storeAttr.Set("id", attr.Id)
		storeAttr.SetFK("model", m.Id)
		tx.h = tx.h.Insert(storeAttr)
	}

	for rKey, rel := range m.Relationships {
		storeRel := model.RecordForModel(RelationshipModel)
		storeRel.Set("name", rKey)
		storeRel.Set("targetModel", rel.TargetModel)
		storeRel.Set("targetRel", rel.TargetRel)
		storeRel.Set("relType", int64(rel.RelType))
		storeRel.Set("id", rel.Id)
		storeRel.SetFK("model", m.Id)
		tx.h = tx.h.Insert(storeRel)
	}
}

func (tx *holdTx) MakeRecord(modelName string) model.Record {
	modelName = strings.ToLower(modelName)
	m, _ := tx.GetModel(modelName)
	rec := model.RecordForModel(m)
	return rec
}

func (tx *holdTx) Resolve(ir *model.IncludeResult, i Inclusion) {
	rec := ir.Record
	id := ir.Record.Id()
	var m q.Matcher
	rel := i.Relationship
	backRel := getBackref(tx, rel)

	switch rel.RelType {
	case model.HasOne:
		// FK on the other side
		targetFK := model.JsonKeyToRelFieldName(rel.TargetRel)
		m = q.Eq(targetFK, id)
		mi := tx.h.IterMatches(rel.TargetModel, m)
		var hits []model.Record
		for val, ok := mi.Next(); ok; val, ok = mi.Next() {
			hits = append(hits, val)
		}
		if len(hits) != 1 {
			panic("Wrong number of hits on hasOne")
		}
		ir.SingleIncludes[backRel.TargetRel] = hits[0]
	case model.BelongsTo:
		// FK on this side
		thisFK := rec.GetFK(backRel.TargetRel)
		m = q.Eq("Id", thisFK)
		mi := tx.h.IterMatches(rel.TargetModel, m)
		var hits []model.Record
		for val, ok := mi.Next(); ok; val, ok = mi.Next() {
			hits = append(hits, val)
		}
		if len(hits) != 1 {
			panic("Wrong number of hits on belongTO")
		}
		ir.SingleIncludes[backRel.TargetRel] = hits[0]
	case model.HasMany:
		// FK on the other side
		targetFK := model.JsonKeyToRelFieldName(rel.TargetRel)
		m = q.Eq(targetFK, id)
		mi := tx.h.IterMatches(rel.TargetModel, m)
		hits := []model.Record{}
		for val, ok := mi.Next(); ok; val, ok = mi.Next() {
			hits = append(hits, val)
		}
		ir.MultiIncludes[backRel.TargetRel] = hits
	case model.HasManyAndBelongsToMany:
		panic("Not implemented")
	}
}

func (tx *holdTx) Commit() {
	tx.ensureWrite()
	tx.db.Lock()
	tx.db.h = tx.h
	tx.db.Unlock()
}
