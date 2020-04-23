package db

import (
	"awans.org/aft/er"
	"awans.org/aft/er/q"
	"awans.org/aft/internal/model"
	"fmt"
	"github.com/google/uuid"
	"github.com/ompluscator/dynamic-struct"
	"reflect"
	"strings"
)

var ModelModel = model.Model{
	Type: "model",
	Id:   uuid.MustParse("872f8c55-9c12-43d1-b3f6-f7a02d937314"),
	Name: "model",
	Attributes: map[string]model.Attribute{
		"name": model.Attribute{
			Id:       uuid.MustParse("d62d3c3a-0228-4131-98f5-2d49a2e3676a"),
			Type:     "attribute",
			AttrType: model.String,
		},
	},
	Relationships: map[string]model.Relationship{
		"attributes": model.Relationship{
			Id:          uuid.MustParse("3271d6a5-0004-4752-81b8-b00142fd59bf"),
			Type:        "relationship",
			TargetModel: "attribute",
			TargetRel:   "model",
			RelType:     model.HasMany,
		},
		"relationships": model.Relationship{
			Id:          uuid.MustParse("806334bf-98ce-4c08-87f4-5d9bed4f6d60"),
			Type:        "relationship",
			TargetModel: "relationship",
			TargetRel:   "model",
			RelType:     model.HasMany,
		},
	},
}

var AttributeModel = model.Model{
	Type: "model",
	Id:   uuid.MustParse("14d840f5-344f-4e23-af12-d4caa1ffa848"),
	Name: "attribute",
	Attributes: map[string]model.Attribute{
		"name": model.Attribute{
			Id:       uuid.MustParse("51605ada-5326-4cfd-9f31-f10bc4dfbf03"),
			Type:     "attribute",
			AttrType: model.String,
		},
		"attrType": model.Attribute{
			Id:       uuid.MustParse("c29a6558-7676-40a8-be00-e0933342efd7"),
			Type:     "attribute",
			AttrType: model.Int,
		},
	},
	Relationships: map[string]model.Relationship{
		"model": model.Relationship{
			Id:          uuid.MustParse("2dbba7d9-3fb0-4905-89f0-d3576e850c05"),
			Type:        "relationship",
			TargetModel: "model",
			TargetRel:   "attributes",
			RelType:     model.BelongsTo,
		},
	},
}

var RelationshipModel = model.Model{
	Type: "model",
	Id:   uuid.MustParse("90be6901-60a0-4eca-893e-232dc57b0bc1"),
	Name: "relationship",
	Attributes: map[string]model.Attribute{
		"name": model.Attribute{
			Id:       uuid.MustParse("7183180e-e13a-4106-844a-04159a8b637c"),
			Type:     "attribute",
			AttrType: model.String,
		},
		"targetModel": model.Attribute{
			Id:       uuid.MustParse("b45e487a-9ed7-4f7d-a760-28691b58e93f"),
			Type:     "attribute",
			AttrType: model.String,
		},
		"targetRel": model.Attribute{
			Id:       uuid.MustParse("3e649bba-b5ab-4ee2-a4ef-3da0eed541da"),
			Type:     "attribute",
			AttrType: model.String,
		},
		"relType": model.Attribute{
			Id:       uuid.MustParse("3c0b2893-a074-4fd7-931e-9a0e45956b08"),
			Type:     "attribute",
			AttrType: model.Int,
		},
	},
	Relationships: map[string]model.Relationship{
		"model": model.Relationship{
			Id:          uuid.MustParse("46962d64-efea-4cde-bad3-bd0170d0866c"),
			Type:        "relationship",
			TargetModel: "model",
			TargetRel:   "relationships",
			RelType:     model.BelongsTo,
		},
	},
}

func New() DB {
	return DB{h: er.New()}
}

func (db DB) AddMetaModel() {
	db.SaveModel(ModelModel)
	db.SaveModel(AttributeModel)
	db.SaveModel(RelationshipModel)
}

type DB struct {
	h *er.Hold
}

func (db DB) GetModel(modelName string) model.Model {
	modelName = strings.ToLower(modelName)
	storeModel, err := db.h.FindOne("model", q.Eq("Name", modelName))
	if err != nil {
		panic(err)
	}
	smReader := dynamicstruct.NewReader(storeModel)

	m := model.Model{
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
			AttrType: model.AttrType(saReader.GetField("Attrtype").Interface().(int)),
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
			RelType:     model.RelType(srReader.GetField("Reltype").Interface().(int)),
			TargetModel: srReader.GetField("Targetmodel").Interface().(string),
			TargetRel:   srReader.GetField("Targetrel").Interface().(string),
		}
		name := srReader.GetField("Name").Interface().(string)
		rels[name] = rel
	}
	m.Relationships = rels

	return m
}

// Manual serialization required for bootstrapping
func (db DB) SaveModel(m model.Model) {
	storeModel := model.StructForModel(ModelModel).New()
	ModelModel.Attributes["name"].SetField("name", m.Name, storeModel)
	model.SystemAttrs["id"].SetField("id", m.Id, storeModel)
	model.SystemAttrs["type"].SetField("type", ModelModel.Name, storeModel)
	db.h.Insert(storeModel)

	for aKey, attr := range m.Attributes {
		storeAttr := model.StructForModel(AttributeModel).New()
		AttributeModel.Attributes["name"].SetField("name", aKey, storeAttr)
		AttributeModel.Attributes["attrType"].SetField("attrType", int(attr.AttrType), storeAttr)
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
		fmt.Printf("before set: %v - %v\n", storeRel, rel.RelType)
		RelationshipModel.Attributes["relType"].SetField("relType", int(rel.RelType), storeRel)
		fmt.Printf("after set: %v\n", storeRel)

		model.SystemAttrs["id"].SetField("id", rel.Id, storeRel)
		model.SystemAttrs["type"].SetField("type", RelationshipModel.Name, storeRel)
		setFK(storeRel, "model", m.Id)
		db.h.Insert(storeRel)
	}
}

func (db DB) MakeStruct(modelName string) interface{} {
	modelName = strings.ToLower(modelName)
	m := db.GetModel(modelName)
	st := model.StructForModel(m).New()
	field := reflect.ValueOf(st).Elem().FieldByName("Type")
	field.SetString(modelName)
	return st
}
