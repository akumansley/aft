package auth

import (
	"awans.org/aft/internal/bizdatatypes"
	"awans.org/aft/internal/db"
)

var UserModel = db.Model{
	ID:   db.MakeModelID("e52f8264-7b95-4a3a-bf76-a23b2229d65a"),
	Name: "user",
	Attributes: []db.Attribute{
		db.Attribute{
			Name:     "email",
			ID:       db.MakeID("236e800d-c39d-4ef3-94e6-5e1f0fc38e62"),
			Datatype: bizdatatypes.EmailAddress,
		},
		db.Attribute{
			Name:     "password",
			ID:       db.MakeID("658f314a-4602-44a9-8d19-884bbd3ea267"),
			Datatype: db.String,
		},
	},
}
