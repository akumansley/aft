package auth

import "awans.org/aft/internal/db"

var LoginSystem = RoleL{
	ID_:      db.MakeID("d8b274db-ee94-4bdd-8c93-b10f64138a8d"),
	Name:     "loginSystem",
	Policies: []PolicyL{LoginUserPolicy, LoginAuthKeyPolicy},
}

var LoginUserPolicy = PolicyL{
	ID_:       db.MakeID("96b1f9fc-7f38-4401-a83f-aa763cc3af0f"),
	AllowRead: true,
	For_:      UserModel,
}

var LoginAuthKeyPolicy = PolicyL{
	ID_:       db.MakeID("97ea6c25-3b91-41ae-a1f7-e1ce0e898698"),
	AllowRead: true,
	For_:      AuthKeyModel,
}
