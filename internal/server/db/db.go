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
)

var (
	ErrData         = errors.New("data-error")
	ErrInvalidModel = fmt.Errorf("%w: invalid model", ErrData)
)

func New() DB {
	appDB := holdDB{h: er.New()}
	appDB.AddMetaModel()
	return appDB
}

func (db holdDB) AddMetaModel() {
	db.SaveModel(ModelModel)
	db.SaveModel(AttributeModel)
	db.SaveModel(RelationshipModel)
}

type DB interface {
	GetModel(string) (model.Model, error)
	SaveModel(model.Model)
	MakeStruct(string) interface{}
	Insert(interface{})
	Connect(from, to interface{}, fromRel model.Relationship)
	Resolve(interface{}, Inclusion)
	FindOne(string, UniqueQuery) (interface{}, error)
	FindMany(string, Query) []interface{}
}

type holdDB struct {
	h *er.Hold
}

func (db holdDB) FindOne(modelName string, uq UniqueQuery) (st interface{}, err error) {
	st, err = db.h.FindOne(modelName, q.Eq(uq.Key, uq.Val))
	return
}

func (db holdDB) Insert(st interface{}) {
	db.h.Insert(st)
}

// TODO hack -- remove this and rewrite with Relationship containing the name
func getBackref(db DB, rel model.Relationship) model.Relationship {
	m, _ := db.GetModel(rel.TargetModel)
	return m.Relationships[rel.TargetRel]
}

func (db holdDB) Connect(from, to interface{}, fromRel model.Relationship) {
	toRel := getBackref(db, fromRel)
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
	db.h.Insert(from)
	db.h.Insert(to)
}

func (db holdDB) GetModel(modelName string) (m model.Model, err error) {
	modelName = strings.ToLower(modelName)
	storeModel, err := db.h.FindOne("model", q.Eq("Name", modelName))
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
	ami := db.h.IterMatches("attribute", q.Eq("ModelId", m.Id))
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
	rmi := db.h.IterMatches("relationship", q.Eq("ModelId", m.Id))
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
func (db holdDB) SaveModel(m model.Model) {
	storeModel := model.StructForModel(ModelModel).New()
	ModelModel.Attributes["name"].SetField("name", m.Name, storeModel)
	model.SystemAttrs["id"].SetField("id", m.Id, storeModel)
	model.SystemAttrs["type"].SetField("type", ModelModel.Name, storeModel)
	db.h.Insert(storeModel)

	for aKey, attr := range m.Attributes {
		storeAttr := model.StructForModel(AttributeModel).New()
		AttributeModel.Attributes["name"].SetField("name", aKey, storeAttr)
		AttributeModel.Attributes["attrType"].SetField("attrType", int64(attr.AttrType), storeAttr)
		model.SystemAttrs["id"].SetField("id", attr.Id, storeAttr)
		model.SystemAttrs["type"].SetField("type", AttributeModel.Name, storeAttr)
		setFK(storeAttr, "model", m.Id)
		db.h.Insert(storeAttr)
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
		db.h.Insert(storeRel)
	}
}

func (db holdDB) MakeStruct(modelName string) interface{} {
	modelName = strings.ToLower(modelName)
	m, _ := db.GetModel(modelName)
	st := model.StructForModel(m).New()
	field := reflect.ValueOf(st).Elem().FieldByName("Type")
	field.SetString(modelName)
	return st
}
