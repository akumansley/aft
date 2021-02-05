package auth

import (
	"fmt"

	"awans.org/aft/internal/db"
)

var UserModel = db.MakeModel(
	db.MakeID("e52f8264-7b95-4a3a-bf76-a23b2229d65a"),
	"user",
	[]db.AttributeL{
		db.MakeConcreteAttribute(
			db.MakeID("236e800d-c39d-4ef3-94e6-5e1f0fc38e62"),
			"email",
			EmailAddress,
		),
		db.MakeConcreteAttribute(
			db.MakeID("658f314a-4602-44a9-8d19-884bbd3ea267"),
			"password",
			Password,
		),
	},
	[]db.RelationshipL{UserRole},
	[]db.ConcreteInterfaceL{},
)

var UserRole = db.MakeConcreteRelationship(
	db.MakeID("e5eea00e-7030-4e6c-85f3-ae8657f365a4"),
	"role",
	false,
	RoleModel,
)

type user struct {
	rec db.Record
	tx  db.Tx
}

func (u *user) ID() db.ID {
	return u.rec.ID()
}

func (u *user) Password() []byte {
	return u.rec.MustGet("password").([]byte)
}

func MakeUser(id db.ID, email, password string, role RoleL) UserL {
	hashedPW, err := hashPassword(id, password)
	if err != nil {
		panic(err)
	}
	return UserL{
		ID_:      id,
		Email:    email,
		Password: hashedPW,
		Role:     role,
	}

}

type UserL struct {
	ID_      db.ID  `record:"id"`
	Email    string `record:"email"`
	Password []byte `record:"password"`
	Role     RoleL
}

func (lit UserL) ID() db.ID {
	return lit.ID_
}

func (lit UserL) InterfaceID() db.ID {
	return UserModel.ID()
}

func (lit UserL) InterfaceName() string {
	return UserModel.Name_
}

func (lit UserL) String() string {
	return fmt.Sprintf("user{%v, %v}", lit.ID_, lit.Email)

}

func (lit UserL) MarshalDB(b *db.Builder) (recs []db.Record, links []db.Link) {
	rec := db.MarshalRecord(b, lit)
	r := lit.Role
	links = append(links, db.Link{From: lit, To: r, Rel: UserRole})
	recs = append(recs, rec)
	return
}

func RoleForUser(tx db.Tx, user db.Record) (db.Record, error) {
	roles := tx.Ref(RoleModel.ID())
	users := tx.Ref(UserModel.ID())
	roleUsers, err := tx.Schema().GetRelationshipByID(RoleUsers.ID())
	if err != nil {
		return nil, err
	}

	q := tx.Query(roles,
		db.Join(users, roles.Rel(roleUsers)),
		db.Aggregate(users, db.Some),
		db.Filter(users, db.EqID(user.ID())),
	)
	roleRec, err := q.OneRecord()

	return roleRec, err
}
