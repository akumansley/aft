package auth

import "awans.org/aft/internal/db"

var Public = RoleL{
	ID_:      db.MakeID("4fb58a4b-ee77-48a2-b63c-164fd7e1f03a"),
	Name:     "public",
	Policies: []PolicyL{},
}
