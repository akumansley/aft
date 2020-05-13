package db

import (
	"awans.org/aft/internal/model"
	"github.com/google/uuid"
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
			AttrType: model.Enum,
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
