package operations

import (
	"awans.org/aft/internal/api"
	"awans.org/aft/internal/db"
	"github.com/stretchr/testify/assert"
	"testing"
)

type DeleteCase struct {
	count     int
	modelName string
}

func TestDeleteApply(t *testing.T) {
	appDB := db.NewTest()
	db.AddSampleModels(appDB)
	tx := appDB.NewRWTx()
	u := api.MakeRecord(tx, "user", `{ 
					"id": "c954af23-a6d7-4930-89e2-87f8c818cc15",
					"type": "user",
					"firstName":"Andrew",
					"lastName":"Wansley",
					"emailAddress":"andrew.wansley@gmail.com",
					"age": 32}`)
	p := api.MakeRecord(tx, "profile", `{
		"id":"d771baa1-3acf-485b-8b47-0d6474a36dee",
		"type":"profile",
		"text": "My bio.."}`)

	up, _ := u.Interface().RelationshipByName("profile")

	var deleteTests = []struct {
		op     DeleteOperation
		output DeleteCase
	}{
		// Simple Delete
		{
			op: DeleteOperation{
				ModelID: db.User.ID(),
				FindManyArgs: FindManyArgs{
					Where: Where{
						FieldCriteria: []FieldCriterion{
							FieldCriterion{
								Key: "firstName",
								Val: "Andrew",
							},
						},
					},
				},
				Nested: []NestedOperation{},
			},
			output: DeleteCase{
				count:     0,
				modelName: "user",
			},
		},
		// Simple Delete fails filter
		{
			op: DeleteOperation{
				ModelID: db.User.ID(),
				FindManyArgs: FindManyArgs{
					Where: Where{
						FieldCriteria: []FieldCriterion{
							FieldCriterion{
								Key: "firstName",
								Val: "Bob",
							},
						},
					},
				},
				Nested: []NestedOperation{},
			},
			output: DeleteCase{
				count:     1,
				modelName: "user",
			},
		},

		// Nested Delete
		{
			op: DeleteOperation{
				ModelID: db.User.ID(),
				FindManyArgs: FindManyArgs{
					Where: Where{
						FieldCriteria: []FieldCriterion{
							FieldCriterion{
								Key: "firstName",
								Val: "Andrew",
							},
						},
					},
				},
				Nested: []NestedOperation{
					NestedDeleteOperation{
						Relationship: up,
					},
				},
			},
			output: DeleteCase{
				count:     0,
				modelName: "profile",
			},
		},

		// Nested Delete with a filter preventing it from deleting
		{
			op: DeleteOperation{
				ModelID: db.User.ID(),
				FindManyArgs: FindManyArgs{
					Where: Where{
						FieldCriteria: []FieldCriterion{
							FieldCriterion{
								Key: "firstName",
								Val: "Andrew",
							},
						},
					},
				},
				Nested: []NestedOperation{
					NestedDeleteOperation{
						Relationship: up,
						Where: Where{
							FieldCriteria: []FieldCriterion{
								FieldCriterion{
									Key: "text",
									Val: "Andrew",
								},
							},
						},
					},
				},
			},
			output: DeleteCase{
				count:     1,
				modelName: "profile",
			},
		},

		// Nested DeleteMany
		{
			op: DeleteOperation{
				ModelID: db.User.ID(),
				FindManyArgs: FindManyArgs{
					Where: Where{
						FieldCriteria: []FieldCriterion{
							FieldCriterion{
								Key: "firstName",
								Val: "Andrew",
							},
						},
					},
				},
				Nested: []NestedOperation{
					NestedDeleteManyOperation{
						Relationship: up,
					},
				},
			},
			output: DeleteCase{
				count:     0,
				modelName: "profile",
			},
		},

		// Nested DeleteMany with a filter preventing it from deleting
		{
			op: DeleteOperation{
				ModelID: db.User.ID(),
				FindManyArgs: FindManyArgs{
					Where: Where{
						FieldCriteria: []FieldCriterion{
							FieldCriterion{
								Key: "firstName",
								Val: "Andrew",
							},
						},
					},
				},
				Nested: []NestedOperation{
					NestedDeleteManyOperation{
						Relationship: up,
						Where: Where{
							FieldCriteria: []FieldCriterion{
								FieldCriterion{
									Key: "text",
									Val: "Andrew",
								},
							},
						},
					},
				},
			},
			output: DeleteCase{
				count:     1,
				modelName: "profile",
			},
		},
	}
	for _, testCase := range deleteTests {
		// start each test on a fresh db
		appDB = db.NewTest()
		db.AddSampleModels(appDB)
		tx = appDB.NewRWTx()
		tx.Insert(u)
		tx.Insert(p)
		tx.Connect(u.ID(), p.ID(), db.UserProfile.ID())
		tx.Commit()
		testCase.op.Apply(tx)
		m, _ := tx.Schema().GetModel(testCase.output.modelName)
		mref := tx.Ref(m.ID())
		found := tx.Query(mref).All()
		assert.Equal(t, testCase.output.count, len(found))
	}
}
