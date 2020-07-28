package operations

import (
	"awans.org/aft/internal/db"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"testing"
)

var (
	profileId = uuid.MustParse("2439b6ce-4dce-4430-8a81-3fe8b7a34ba1")
)

func addIncludeTestData(appDB db.DB) {
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
	tx.Connect(u1.ID(), pr.ID(), db.UserProfile.ID())

	tx.Commit()
}

func TestInclude(t *testing.T) {
	appDB := db.NewTest()
	db.AddSampleModels(appDB)
	addIncludeTestData(appDB)
	tx := appDB.NewTx()
	var includeTests = []struct {
		operation FindManyOperation
		output    []uuid.UUID
	}{
		// Simple Include
		{
			operation: FindManyOperation{
				ModelID: db.User.ID(),
				Where:   Where{},
				Include: Include{
					[]Inclusion{
						Inclusion{
							Relationship:   db.UserProfile,
							NestedFindMany: NestedFindManyOperation{},
						},
					},
				},
			},
			output: []uuid.UUID{profileId},
		},

		// Nested Include
		//		{
		//			operation: FindManyOperation{
		//				ModelID: db.Profile.ID(),
		//				Where:   Where{},
		//				Include: Include{
		//					[]Inclusion{
		//						Inclusion{
		//							Relationship: db.ProfileUser,
		//							NestedFindMany: NestedFindManyOperation{
		//								Where: Where{},
		//								Include: Include{
		//									[]Inclusion{
		//										Inclusion{
		//											Relationship: db.UserPosts,
		//											NestedFindMany: NestedFindManyOperation{
		//												Where: Where{},
		//											},
		//										},
		//									},
		//								},
		//							},
		//						},
		//					},
		//				},
		//			},
		//			output: []uuid.UUID{userId1, postId1},
		//		},
	}
	for _, testCase := range includeTests {
		records, _ := testCase.operation.Apply(tx)
		for _, v := range records {
			if v.ToOne["profile"] != nil {
				assert.ElementsMatch(t, testCase.output[0], v.ToOne["profile"].Record.ID())
			}
			if v.ToOne["user"] != nil {
				out := testCase.output[1].String()
				post1 := v.ToOne["user"].ToMany["posts"][0].Record.ID().String()
				post2 := v.ToOne["user"].ToMany["posts"][1].Record.ID().String()
				assert.True(t, out == post1 || out == post2)
			}
		}
	}
}
