package operations

import (
	"testing"

	"awans.org/aft/internal/api"
	"awans.org/aft/internal/db"
	"github.com/stretchr/testify/assert"
)

type SelectCase struct {
	count        int
	relationship string
}

func TestSelect(t *testing.T) {
	appDB := db.NewTest()
	db.AddSampleModels(appDB)
	tx := appDB.NewRWTx()
	u1, err := tx.MakeRecord(db.User.ID())
	if err != nil {
		panic(err)
	}
	u1.Set("id", userId1)
	u1.Set("firstName", "Gid")
	u1.Set("age", int64(4))

	u2, err := tx.MakeRecord(db.User.ID())
	if err != nil {
		panic(err)
	}
	u2.Set("id", userId2)
	u2.Set("firstName", "Chase")
	u2.Set("age", int64(5))

	u3, err := tx.MakeRecord(db.User.ID())
	if err != nil {
		panic(err)
	}
	u3.Set("id", userId3)
	u3.Set("firstName", "Tom")
	u3.Set("age", int64(6))

	p1, err := tx.MakeRecord(db.Post.ID())
	if err != nil {
		panic(err)
	}
	p1.Set("id", postId1)
	p1.Set("text", "hello")

	p2, err := tx.MakeRecord(db.Post.ID())
	if err != nil {
		panic(err)
	}
	p2.Set("id", postId2)
	p2.Set("text", "goodbye")

	pr, err := tx.MakeRecord(db.Profile.ID())
	if err != nil {
		panic(err)
	}
	pr.Set("id", profileId)
	pr.Set("text", "cool")

	tx.Insert(u1)
	tx.Insert(u2)
	tx.Insert(u3)
	tx.Insert(p1)
	tx.Insert(p2)
	tx.Insert(pr)
	tx.Connect(u1.ID(), p1.ID(), db.UserPosts.ID())
	tx.Connect(u1.ID(), p2.ID(), db.UserPosts.ID())
	tx.Connect(pr.ID(), u1.ID(), db.ProfileUser.ID())
	tx.Connect(u1.ID(), pr.ID(), db.UserProfile.ID())

	tx.Commit()
	upr, _ := tx.Schema().GetRelationshipByID(db.UserProfile.ID())

	fields := make(api.Set)
	fields["firstName"] = api.Void{}
	var selectTests = []struct {
		operation FindManyOperation
		output    SelectCase
	}{
		// Simple Select
		{
			operation: FindManyOperation{
				ModelID: db.User.ID(),
				FindArgs: FindArgs{
					Where: Where{},
					Select: Select{
						true,
						fields,
						[]Selection{
							Selection{
								Relationship:   upr,
								NestedFindMany: FindArgs{},
							},
						},
					},
				},
			},
			output: SelectCase{
				count:        1,
				relationship: "profile",
			},
		},
	}
	for _, testCase := range selectTests {
		records, _ := testCase.operation.Apply(tx)
		for _, v := range records {
			assert.NotContains(t, v.String(), "lastName")
		}
	}
}
