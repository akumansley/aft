package auth

import (
	"awans.org/aft/internal/db"
)

var AuthKeyModel = db.Model{
	ID:     db.MakeModelID("0285e736-6a8b-4c47-852e-e73f12eb94f4"),
	Name:   "authKey",
	System: true,
	Attributes: []db.Attribute{
		db.Attribute{
			Name:     "key",
			ID:       db.MakeID("84ece8dd-076b-493b-8199-f1ea2ca5acb7"),
			Datatype: db.String,
		},
		db.Attribute{
			Name:     "active",
			ID:       db.MakeID("ca89e9a6-b613-4c8e-9154-d4c6d3334c9a"),
			Datatype: db.Bool,
		},
	},
}
