package auth

import (
	"testing"

	"awans.org/aft/internal/db"
)

var giftListItem = db.MakeModel(
	db.MakeID("e76b3254-d0fc-453b-8269-2bc05b4c83b9"),
	"giftListItem",
	[]db.AttributeL{},
	[]db.RelationshipL{},
	[]db.ConcreteInterfaceL{},
)

var comment = db.MakeModel(
	db.MakeID("6db17669-7d8b-4c26-befa-f7594621aa33"),
	"comment",
	[]db.AttributeL{},
	[]db.RelationshipL{},
	[]db.ConcreteInterfaceL{},
)

var user1 = MakeUser(
	db.MakeID("ff057685-8f22-4c28-b150-3d5efa1b3d3a"),
	"user1@gmail.com",
	"coolpass",
	userRole,
)

var user2 = MakeUser(
	db.MakeID("f69d4007-0eca-44ad-8666-48ccdab02717"),
	"user2@gmail.com",
	"coolpass",
	userRole,
)

var adminUser = MakeUser(
	db.MakeID("919f9ab2-687a-4957-8e12-359c28df2b16"),
	"admin@gmail.com",
	"coolpass",
	adminRole,
)

var userRole = RoleL{
	ID_:      db.MakeID("dc70fec2-8ff8-4983-9c7b-34c798a88f8a"),
	Name:     "user",
	Policies: []PolicyL{signedinPolicy},
}

var adminRole = RoleL{
	ID_:      db.MakeID("4aaa6f9f-fd6d-4d80-97a1-b0aa2ffdce52"),
	Name:     "admin",
	Policies: []PolicyL{adminPolicy},
}

var ownerPolicy = PolicyL{
	ID_:       db.MakeID("f91739a8-a727-48b7-a951-49e3cbbfe37e"),
	AllowRead: true,
	ReadWhere: ``,
	For_:      UserModel,
}

var anyoneButOwnerPolicy = PolicyL{
	ID_:       db.MakeID("09d881bc-5246-4197-847b-037b55c2e5b0"),
	AllowRead: true,
	ReadWhere: ``,
	For_:      UserModel,
}

func TestPolicy(t *testing.T) {
}
