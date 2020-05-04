package db

import (
	"awans.org/aft/er"
	"awans.org/aft/er/q"
	"awans.org/aft/internal/model"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/ompluscator/dynamic-struct"
	"reflect"
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
	MakeStruct(string) interface{}
	Resolve(interface{}, Inclusion)
	FindOne(string, UniqueQuery) (interface{}, error)
	FindMany(string, Query) []interface{}
}

type RWTx interface {
	GetModel(string) (model.Model, error)
	MakeStruct(string) interface{}
	Resolve(interface{}, Inclusion)
	FindOne(string, UniqueQuery) (interface{}, error)
	FindMany(string, Query) []interface{}

	SaveModel(model.Model)
	Insert(interface{})
	Connect(from, to interface{}, fromRel model.Relationship)
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

func (tx *holdTx) FindOne(modelName string, uq UniqueQuery) (st interface{}, err error) {
	st, err = tx.h.FindOne(modelName, q.Eq(uq.Key, uq.Val))
	return
}

func (tx *holdTx) Insert(st interface{}) {
	tx.ensureWrite()
	tx.h = tx.h.Insert(st)
}

// TODO hack -- remove this and rewrite with Relationship containing the name
func getBackref(tx Tx, rel model.Relationship) model.Relationship {
	m, _ := tx.GetModel(rel.TargetModel)
	return m.Relationships[rel.TargetRel]
}

func (tx *holdTx) Connect(from, to interface{}, fromRel model.Relationship) {
	tx.ensureWrite()
	toRel := getBackref(tx, fromRel)
	if fromRel.RelType == model.BelongsTo && (toRel.RelType == model.HasOne || toRel.RelType == model.HasMany) {
		// FK from
		setFK(from, toRel.TargetRel, getId(to))
	} else if toRel.RelType == model.BelongsTo && (fromRel.RelType == model.HasOne || fromRel.RelType == model.HasMany) {
		// FK to
		setFK(to, fromRel.TargetRel, getId(from))
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
	storeModel, err := tx.h.FindOne("model", q.Eq("Name", modelName))
	if err != nil {
		return m, fmt.Errorf("%w: %v", ErrInvalidModel, modelName)
	}
	smReader := dynamicstruct.NewReader(storeModel)

	m = model.Model{
		Type: smReader.GetField("Type").Interface().(string),
		Id:   smReader.GetField("Id").Interface().(uuid.UUID),
		Name: smReader.GetField("Name").Interface().(string),
	}

	attrs := make(map[string]model.Attribute)

	// make ModelId a dynamic key
	ami := tx.h.IterMatches("attribute", q.Eq("ModelId", m.Id))
	for storeAttr, ok := ami.Next(); ok; storeAttr, ok = ami.Next() {
		saReader := dynamicstruct.NewReader(storeAttr)
		attr := model.Attribute{
			AttrType: model.AttrType(saReader.GetField("Attrtype").Interface().(int64)),
			Type:     saReader.GetField("Type").Interface().(string),
			Id:       saReader.GetField("Id").Interface().(uuid.UUID),
		}
		name := saReader.GetField("Name").Interface().(string)
		attrs[name] = attr
	}
	m.Attributes = attrs

	rels := make(map[string]model.Relationship)

	// make ModelId a dynamic key
	rmi := tx.h.IterMatches("relationship", q.Eq("ModelId", m.Id))
	for storeRel, ok := rmi.Next(); ok; storeRel, ok = rmi.Next() {
		srReader := dynamicstruct.NewReader(storeRel)
		rel := model.Relationship{
			Type:        srReader.GetField("Type").Interface().(string),
			Id:          srReader.GetField("Id").Interface().(uuid.UUID),
			RelType:     model.RelType(srReader.GetField("Reltype").Interface().(int64)),
			TargetModel: srReader.GetField("Targetmodel").Interface().(string),
			TargetRel:   srReader.GetField("Targetrel").Interface().(string),
		}
		name := srReader.GetField("Name").Interface().(string)
		rels[name] = rel
	}
	m.Relationships = rels

	return m, nil
}

// Manual serialization required for bootstrapping
func (tx *holdTx) SaveModel(m model.Model) {
	tx.ensureWrite()
	storeModel := model.StructForModel(ModelModel).New()
	ModelModel.Attributes["name"].SetField("name", m.Name, storeModel)
	model.SystemAttrs["id"].SetField("id", m.Id, storeModel)
	model.SystemAttrs["type"].SetField("type", ModelModel.Name, storeModel)
	tx.h = tx.h.Insert(storeModel)

	for aKey, attr := range m.Attributes {
		storeAttr := model.StructForModel(AttributeModel).New()
		AttributeModel.Attributes["name"].SetField("name", aKey, storeAttr)
		AttributeModel.Attributes["attrType"].SetField("attrType", int64(attr.AttrType), storeAttr)
		model.SystemAttrs["id"].SetField("id", attr.Id, storeAttr)
		model.SystemAttrs["type"].SetField("type", AttributeModel.Name, storeAttr)
		setFK(storeAttr, "model", m.Id)
		tx.h = tx.h.Insert(storeAttr)
	}

	for rKey, rel := range m.Relationships {
		storeRel := model.StructForModel(RelationshipModel).New()
		RelationshipModel.Attributes["name"].SetField("name", rKey, storeRel)
		RelationshipModel.Attributes["targetModel"].SetField("targetModel", rel.TargetModel, storeRel)
		RelationshipModel.Attributes["targetRel"].SetField("targetRel", rel.TargetRel, storeRel)
		RelationshipModel.Attributes["relType"].SetField("relType", int64(rel.RelType), storeRel)

		model.SystemAttrs["id"].SetField("id", rel.Id, storeRel)
		model.SystemAttrs["type"].SetField("type", RelationshipModel.Name, storeRel)
		setFK(storeRel, "model", m.Id)
		tx.h = tx.h.Insert(storeRel)
	}
}

func (tx *holdTx) MakeStruct(modelName string) interface{} {
	modelName = strings.ToLower(modelName)
	m, _ := tx.GetModel(modelName)
	st := model.StructForModel(m).New()
	field := reflect.ValueOf(st).Elem().FieldByName("Type")
	field.SetString(modelName)
	return st
}

func (tx *holdTx) Resolve(st interface{}, i Inclusion) {
	id := getId(st)
	var m q.Matcher
	rel := i.Relationship
	backRel := getBackref(tx, rel)
	var related interface{}
	switch rel.RelType {
	case model.HasOne:
		// FK on the other side
		targetFK := model.JsonKeyToRelFieldName(rel.TargetRel)
		m = q.Eq(targetFK, id)
		mi := tx.h.IterMatches(rel.TargetModel, m)
		var hits []interface{}
		for val, ok := mi.Next(); ok; val, ok = mi.Next() {
			hits = append(hits, val)
		}
		if len(hits) != 1 {
			panic("Wrong number of hits on hasOne")
		}
		related = hits[0]
	case model.BelongsTo:
		// FK on this side
		thisFK := getFK(st, backRel.TargetRel)
		m = q.Eq("Id", thisFK)
		mi := tx.h.IterMatches(rel.TargetModel, m)
		var hits []interface{}
		for val, ok := mi.Next(); ok; val, ok = mi.Next() {
			hits = append(hits, val)
		}
		if len(hits) != 1 {
			panic("Wrong number of hits on belongTO")
		}
		related = hits[0]
	case model.HasMany:
		// FK on the other side
		targetFK := model.JsonKeyToRelFieldName(rel.TargetRel)
		m = q.Eq(targetFK, id)
		mi := tx.h.IterMatches(rel.TargetModel, m)
		hits := []interface{}{}
		for val, ok := mi.Next(); ok; val, ok = mi.Next() {
			hits = append(hits, val)
		}
		related = hits
	case model.HasManyAndBelongsToMany:
		panic("Not implemented")
	}
	setRelated(st, backRel.TargetRel, related)

}

func setRelated(st interface{}, key string, val interface{}) {
	fieldName := model.JsonKeyToFieldName(key)
	field := reflect.ValueOf(st).Elem().FieldByName(fieldName)
	v := reflect.ValueOf(&val)
	field.Set(v)
}

func (tx *holdTx) Commit() {
	tx.ensureWrite()
	tx.db.Lock()
	tx.db.h = tx.h
	tx.db.Unlock()
}
