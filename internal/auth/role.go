package auth

import (
	"awans.org/aft/internal/db"
)

// TODO: relationship between user and roles
var RoleModel = db.MakeModel(
	db.MakeID("bf17994e-7ef1-459f-9b82-069016686081"),
	"role",
	[]db.AttributeL{
		db.MakeConcreteAttribute(
			db.MakeID("6dc3ec26-3125-4e54-b9b0-f5ccad10c4af"),
			"name",
			db.String,
		),
	},
	[]db.RelationshipL{},
	[]db.ConcreteInterfaceL{},
)

func init() {
	RoleModel.Relationships_ = []db.RelationshipL{
		RoleUsers,
	}
}

var RoleUsers = db.MakeReverseRelationship(
	db.MakeID("098dd9f8-1337-44b2-bf8d-277e4aafd725"),
	"users",
	UserRoles,
)
