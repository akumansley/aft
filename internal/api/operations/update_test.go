package operations

import (
	"awans.org/aft/internal/api"
	"awans.org/aft/internal/db"
	"github.com/stretchr/testify/assert"
	"testing"
)

type UpdateCase struct {
	count int
}

func TestUpdateApply(t *testing.T) {
	appDB := db.NewTest()
	db.AddSampleModels(appDB)
	tx := appDB.NewRWTx()
	u := api.MakeRecord(tx, "user", `{ 
					"type": "user",
					"firstName":"Andrew",
					"lastName":"Wansley",
					"emailAddress":"andrew.wansley@gmail.com",
					"age": 32}`)
	p := api.MakeRecord(tx, "profile", `{
		"type":"profile",
		"text": "My bio.."}`)
	up, _ := u.Interface().RelationshipByName("profile")

	var updateTests = []struct {
		op     UpdateOperation
		output UpdateCase
	}{
		// Simple update
		{
			op: UpdateOperation{
				ModelID: db.User.ID(),
				Where: Where{
					FieldCriteria: []FieldCriterion{
						FieldCriterion{
							Key: "firstName",
							Val: "Andrew",
						},
					},
				},
				Data: map[string]interface{}{
					"firstName": "bob",
				},
				Nested: []NestedOperation{},
			},
			output: UpdateCase{
				count: 0,
			},
		},

		// Nested update
		{
			op: UpdateOperation{
				ModelID: db.User.ID(),
				Where: Where{
					FieldCriteria: []FieldCriterion{
						FieldCriterion{
							Key: "firstName",
							Val: "Andrew",
						},
					},
				},
				Data: map[string]interface{}{
					"firstName": "bob",
				},
				Nested: []NestedOperation{
					NestedUpdateOperation{
						Relationship: up,
						Data: map[string]interface{}{
							"text": "cool",
						},
					},
				},
			},
			output: UpdateCase{
				count: 1,
			},
		},

		// Nested deleteMany
		{
			op: UpdateOperation{
				ModelID: db.User.ID(),
				Where: Where{
					FieldCriteria: []FieldCriterion{
						FieldCriterion{
							Key: "firstName",
							Val: "Andrew",
						},
					},
				},
				Data: map[string]interface{}{
					"firstName": "bob",
				},
				Nested: []NestedOperation{
					NestedDeleteManyOperation{
						Relationship: up,
					},
				},
			},
			output: UpdateCase{
				count: 0,
			},
		},

		// Nested delete
		{
			op: UpdateOperation{
				ModelID: db.User.ID(),
				Where: Where{
					FieldCriteria: []FieldCriterion{
						FieldCriterion{
							Key: "firstName",
							Val: "Andrew",
						},
					},
				},
				Data: map[string]interface{}{
					"firstName": "bob",
				},
				Nested: []NestedOperation{
					NestedDeleteOperation{
						Relationship: up,
					},
				},
			},
			output: UpdateCase{
				count: 0,
			},
		},

		// Nested upsert
		{
			op: UpdateOperation{
				ModelID: db.User.ID(),
				Where: Where{
					FieldCriteria: []FieldCriterion{
						FieldCriterion{
							Key: "firstName",
							Val: "Andrew",
						},
					},
				},
				Data: map[string]interface{}{
					"firstName": "bob",
				},
				Nested: []NestedOperation{
					NestedUpsertOperation{
						Relationship: up,
						Where: Where{
							FieldCriteria: []FieldCriterion{
								FieldCriterion{
									Key: "text",
									Val: "My bio..",
								},
							},
						},
						Update: map[string]interface{}{"text": "cool"},
					},
				},
			},
			output: UpdateCase{
				count: 1,
			},
		},
	}
	for _, testCase := range updateTests {
		// start each test on a fresh db
		appDB = db.NewTest()
		db.AddSampleModels(appDB)
		tx = appDB.NewRWTx()
		tx.Insert(u)
		tx.Insert(p)
		tx.Connect(u.ID(), p.ID(), db.UserProfile.ID())
		tx.Commit()
		out, _ := testCase.op.Apply(tx)
		assert.Equal(t, "bob", out.Record.MustGet("firstName"))
		assert.Equal(t, "Wansley", out.Record.MustGet("lastName"))
		r := tx.Ref(p.Interface().ID())
		q := tx.Query(r, db.Filter(r, db.Eq("text", "cool")))
		assert.Equal(t, testCase.output.count, len(q.All()))
	}
}
