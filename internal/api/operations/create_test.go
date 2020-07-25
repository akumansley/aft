package operations

import (
	"awans.org/aft/internal/api"
	"awans.org/aft/internal/db"
	"github.com/go-test/deep"
	"testing"
)

type CreateCase struct {
	st        db.Record
	modelName string
}

func TestCreateApply(t *testing.T) {
	appDB := db.NewTest()
	db.AddSampleModels(appDB)
	tx := appDB.NewRWTx()
	u := api.MakeRecord(tx, "user", `{ 
					"id"  : "36d52356-5730-4f45-9305-08b799b93c3b",
					"type": "user",
					"firstName":"Andrew",
					"lastName":"Wansley",
					"emailAddress":"andrew.wansley@gmail.com",
					"age": 32}`)
	p := api.MakeRecord(tx, "profile", `{
		"id"  : "c6a1eb80-b532-4bce-a4c5-577544ea9847",
		"type":"profile",
		"text": "My bio.."}`)

	up, _ := u.Interface().RelationshipByName("profile")
	var createTests = []struct {
		operations []CreateOperation
		output     []CreateCase
	}{
		// Simple Create
		{
			operations: []CreateOperation{
				CreateOperation{
					Record: u,
					Nested: []NestedOperation{},
				},
			},
			output: []CreateCase{
				CreateCase{
					st:        u,
					modelName: "user",
				},
			},
		},
		// Nested Create
		{
			operations: []CreateOperation{
				CreateOperation{
					Record: u,
					Nested: []NestedOperation{
						NestedCreateOperation{
							Relationship: up,
							Record:       p,
							Nested:       []NestedOperation{},
						},
					},
				},
			},
			output: []CreateCase{
				CreateCase{
					st:        u,
					modelName: "user",
				},
				CreateCase{
					st:        p,
					modelName: "profile",
				},
			},
		},
		// Nested connect
		{
			operations: []CreateOperation{
				CreateOperation{
					Record: p,
					Nested: []NestedOperation{},
				},
				CreateOperation{
					Record: u,
					Nested: []NestedOperation{
						NestedConnectOperation{
							Relationship: up,
							Where: Where{
								FieldCriteria: []FieldCriterion{
									FieldCriterion{
										Key: "text",
										Val: "My bio..",
									},
								},
							},
						},
					},
				},
			},
			output: []CreateCase{
				CreateCase{
					st:        u,
					modelName: "user",
				},
				CreateCase{
					st:        p,
					modelName: "profile",
				},
			},
		},
	}
	for _, testCase := range createTests {
		// start each test on a fresh db
		appDB = db.NewTest()
		db.AddSampleModels(appDB)
		tx = appDB.NewRWTx()
		for _, op := range testCase.operations {
			op.Apply(tx)
		}
		for _, CreateCase := range testCase.output {
			m, _ := tx.Schema().GetModel(CreateCase.modelName)
			mref := tx.Ref(m.ID())
			found, _ := tx.Query(mref, db.Filter(mref, db.EqID(CreateCase.st.ID()))).OneRecord()
			if diff := deep.Equal(found, CreateCase.st); diff != nil {
				t.Error(diff)
			}
		}

	}
}
