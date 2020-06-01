package auth

import (
	"awans.org/aft/internal/db"
	"github.com/google/uuid"
)

var AuthKeyModel = db.Model{
	Id:   uuid.MustParse("0285e736-6a8b-4c47-852e-e73f12eb94f4"),
	Name: "authKey",
	Attributes: map[string]db.Attribute{
		"key": db.Attribute{
			Id:       uuid.MustParse("84ece8dd-076b-493b-8199-f1ea2ca5acb7"),
			AttrType: db.String,
		},
		"active": db.Attribute{
			Id:       uuid.MustParse("ca89e9a6-b613-4c8e-9154-d4c6d3334c9a"),
			AttrType: db.Bool,
		},
	},
}
