package operations

import (
	"testing"

	"awans.org/aft/internal/api"
	"awans.org/aft/internal/db"
	"github.com/stretchr/testify/assert"
)

type UpdateTest struct {
	op     UpdateOperation
	output UpdateCase
}

type UpdateCase struct {
	count int
}

const numCases = 6

func getTestCase(ix int, tx db.Tx) UpdateTest {
	user, err := tx.Schema().GetModel("user")
	if err != nil {
		panic(err)
	}
	up, err := user.RelationshipByName("profile")
	if err != nil {
		panic(err)
	}

	var updateTests = []UpdateTest{
		// Simple update
		{
			op: UpdateOperation{
				ModelID: db.User.ID(),
				FindArgs: FindArgs{
					Where: Where{
						FieldCriteria: []FieldCriterion{
							FieldCriterion{
								Key: "firstName",
								Val: "Andrew",
							},
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
				FindArgs: FindArgs{
					Where: Where{
						FieldCriteria: []FieldCriterion{
							FieldCriterion{
								Key: "firstName",
								Val: "Andrew",
							},
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
				FindArgs: FindArgs{
					Where: Where{
						FieldCriteria: []FieldCriterion{
							FieldCriterion{
								Key: "firstName",
								Val: "Andrew",
							},
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

		// Nested set
		{
			op: UpdateOperation{
				ModelID: db.User.ID(),
				FindArgs: FindArgs{
					Where: Where{
						FieldCriteria: []FieldCriterion{
							FieldCriterion{
								Key: "firstName",
								Val: "Andrew",
							},
						},
					},
				},
				Data: map[string]interface{}{
					"firstName": "bob",
				},
				Nested: []NestedOperation{
					NestedSetOperation{
						Where: Where{
							FieldCriteria: []FieldCriterion{
								FieldCriterion{
									Key: "text",
									Val: "My bio..",
								},
							},
						},
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
				FindArgs: FindArgs{
					Where: Where{
						FieldCriteria: []FieldCriterion{
							FieldCriterion{
								Key: "firstName",
								Val: "Andrew",
							},
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
				FindArgs: FindArgs{
					Where: Where{
						FieldCriteria: []FieldCriterion{
							FieldCriterion{
								Key: "firstName",
								Val: "Andrew",
							},
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
						Create: map[string]interface{}{"type": "profile", "text": "awes.."},
					},
				},
			},
			output: UpdateCase{
				count: 1,
			},
		},
	}

	return updateTests[ix]
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

	for i := 0; i < numCases; i++ {
		appDB = db.NewTest()
		db.AddSampleModels(appDB)
		tx = appDB.NewRWTx()

		testCase := getTestCase(i, tx)

		tx.Insert(u)
		tx.Insert(p)
		err := tx.Connect(u.ID(), p.ID(), db.UserProfile.ID())
		if err != nil {
			t.Fatal(err)
		}
		out, err := testCase.op.Apply(tx)
		if err != nil {
			t.Fatal(err)
		}
		if out == nil {
			t.Fatal("nil QR from update")
		}
		assert.Equal(t, "bob", out.Record.MustGet("firstName"))
		assert.Equal(t, "Wansley", out.Record.MustGet("lastName"))

		r := tx.Ref(p.InterfaceID())

		q := tx.Query(r, db.Filter(r, db.Eq("text", "cool")))

		assert.Equal(t, testCase.output.count, len(q.All()))
	}
}
