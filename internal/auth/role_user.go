package auth

import "awans.org/aft/internal/db"

var UserRoleL = RoleL{
	ID_:      db.MakeID("aff95271-e7d1-477e-9325-03aabf2f83f9"),
	Name:     "user",
	Policies: []PolicyL{UserUserPolicy},
}

var UserUserPolicy = PolicyL{
	ID_:         db.MakeID("a85f6350-a613-4333-a455-d6e705fe19e5"),
	AllowRead:   true,
	ReadWhere:   `{"id":"$USER_ID"}`,
	AllowUpdate: true,
	UpdateWhere: `{"id":"$USER_ID"}`,
	AllowCreate: false,
	For_:        UserModel,
}
