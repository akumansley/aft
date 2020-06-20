package auth

import (
	"awans.org/aft/internal/db"
)

var PolicyModel = db.MakeModel(
	db.MakeID("ea5eda03-6780-4a31-8b9b-e5f16a98d8b3"),
	"policy",
	[]db.AttributeL{
		db.MakeConcreteAttribute(
			db.MakeID("55cfda72-c7f2-47aa-85ab-e54b98f1eda0"),
			"text",
			db.String,
		),
		db.MakeConcreteAttribute(
			db.MakeID("7ebfbce0-3280-4067-8cce-c00efa89bb43"),
			"name",
			db.String,
		),
	},
	[]db.RelationshipL{PolicyFor},
	[]db.ConcreteInterfaceL{},
)

var PolicyFor = db.MakeConcreteRelationship(
	db.MakeID("be24d5ca-48f4-4d6f-a550-5b969703f440"),
	"model",
	false,
	db.ModelModel,
)
