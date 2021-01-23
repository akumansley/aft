package auth

import (
	"awans.org/aft/internal/db"
)

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
		RolePolicy,
	}
	PolicyModel.Relationships_ = []db.RelationshipL{
		PolicyRole, PolicyFor,
	}
}

var RoleUsers = db.MakeReverseRelationship(
	db.MakeID("098dd9f8-1337-44b2-bf8d-277e4aafd725"),
	"users",
	UserRole,
)

var RolePolicy = db.MakeConcreteRelationship(
	db.MakeID("fc193452-3c43-4019-b886-d95decc1ce97"),
	"policies",
	true,
	PolicyModel,
)

type RoleL struct {
	ID_      db.ID  `record:"id"`
	Name     string `record:"name"`
	Policies []PolicyL
}

func (lit RoleL) ID() db.ID {
	return lit.ID_
}

func (lit RoleL) InterfaceID() db.ID {
	return RoleModel.ID()
}

func (lit RoleL) MarshalDB(b *db.Builder) (recs []db.Record, links []db.Link) {
	rec := db.MarshalRecord(b, lit)
	for _, p := range lit.Policies {
		links = append(links, db.Link{rec.ID(), p.ID(), RolePolicy})
	}
	recs = append(recs, rec)
	return
}
