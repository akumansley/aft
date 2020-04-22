package db

import (
	"awans.org/aft/er"
	"awans.org/aft/er/q"
	"awans.org/aft/internal/model"
	"github.com/google/uuid"
	"reflect"
	"strings"
)

var ModelModel = model.Model{
	Type: "model",
	Id:   uuid.MustParse("872f8c55-9c12-43d1-b3f6-f7a02d937314"),
	Name: "model",
	Attributes: map[string]model.Attribute{
		"name": model.Attribute{
			Type: model.String,
		},
	},
	Relationships: map[string]model.Relationship{
		"attributes": model.Relationship{
			TargetModel: "attribute",
			TargetRel:   "model",
			Type:        model.HasMany,
		},
		"relationships": model.Relationship{
			TargetModel: "relationship",
			TargetRel:   "model",
			Type:        model.HasMany,
		},
	},
}

var AttributeModel = model.Model{
	Type: "model",
	Id:   uuid.MustParse("14d840f5-344f-4e23-af12-d4caa1ffa848"),
	Name: "attribute",
	Attributes: map[string]model.Attribute{
		"name": model.Attribute{
			Type: model.String,
		},
		"type": model.Attribute{
			Type: model.Int,
		},
	},
	Relationships: map[string]model.Relationship{
		"model": model.Relationship{
			TargetModel: "model",
			TargetRel:   "attributes",
			Type:        model.BelongsTo,
		},
	},
}

var RelationshipModel = model.Model{
	Type: "model",
	Id:   uuid.MustParse("90be6901-60a0-4eca-893e-232dc57b0bc1"),
	Name: "relationship",
	Attributes: map[string]model.Attribute{
		"name": model.Attribute{
			Type: model.String,
		},
		"targetModel": model.Attribute{
			Type: model.String,
		},
		"targetRel": model.Attribute{
			Type: model.String,
		},
		"type": model.Attribute{
			Type: model.Int,
		},
	},
	Relationships: map[string]model.Relationship{
		"model": model.Relationship{
			TargetModel: "model",
			TargetRel:   "relationships",
			Type:        model.BelongsTo,
		},
	},
}

func New() DB {
	return DB{h: er.New()}
}

func (db DB) AddMetaModel() {
	db.h.Insert(ModelModel)
	db.h.Insert(RelationshipModel)
	db.h.Insert(AttributeModel)
}

type DB struct {
	h *er.Hold
}

func (db DB) GetModel(modelName string) model.Model {
	modelName = strings.ToLower(modelName)
	val, err := db.h.FindOne("model", q.Eq("Name", modelName))
	if err != nil {
		panic(err)
	}
	m, ok := val.(model.Model)
	if !ok {
		panic("Not a model")
	}
	return m
}

func (db DB) MakeStruct(modelName string) interface{} {
	modelName = strings.ToLower(modelName)
	if modelName == "model" {
		return model.Model{}
	} else {
		m := db.GetModel(modelName)
		st := model.StructForModel(m).New()
		field := reflect.ValueOf(st).Elem().FieldByName("Type")
		field.SetString(modelName)
		return st
	}
}
