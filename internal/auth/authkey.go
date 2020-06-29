package auth

import (
	"awans.org/aft/internal/db"
)

var AuthKeyModel = db.MakeModel(
	db.MakeID("0285e736-6a8b-4c47-852e-e73f12eb94f4"),
	"authKey",
	[]db.AttributeL{
		db.MakeConcreteAttribute(
			db.MakeID("84ece8dd-076b-493b-8199-f1ea2ca5acb7"),
			"key",
			db.String,
		),
		db.MakeConcreteAttribute(
			db.MakeID("ca89e9a6-b613-4c8e-9154-d4c6d3334c9a"),
			"active",
			db.Bool,
		),
	},
	[]db.RelationshipL{},
	[]db.ConcreteInterfaceL{},
)
