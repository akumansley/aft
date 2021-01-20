package auth

import "awans.org/aft/internal/db"

var Public = RoleL{
	ID_:      db.MakeID("4fb58a4b-ee77-48a2-b63c-164fd7e1f03a"),
	Name:     "public",
	Policies: []PolicyL{},
}

func getPublic(tx db.Tx) db.Record {
	roles := tx.Ref(RoleModel.ID())
	val, err := tx.Query(roles, db.Filter(roles, db.EqID(Public.ID()))).OneRecord()
	if err != nil {
		panic(err)
	}
	return val
}

var System = RoleL{
	ID_:      db.MakeID("d8b274db-ee94-4bdd-8c93-b10f64138a8d"),
	Name:     "system",
	Policies: []PolicyL{SystemUserPolicy},
}

var SystemUserPolicy = PolicyL{
	ID_:       db.MakeID("96b1f9fc-7f38-4401-a83f-aa763cc3af0f"),
	AllowRead: true,
	For_:      UserModel,
}
