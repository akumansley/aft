package operations

import (
	"awans.org/aft/internal/db"
	"github.com/stretchr/testify/assert"
	"testing"
)

type CreateCase struct {
	st        map[string]interface{}
	modelName string
}

func TestCreateApply(t *testing.T) {
	appDB := db.NewTest()
	db.AddSampleModels(appDB)
	tx := appDB.NewRWTx()
	u := map[string]interface{}{
		"firstName":    "Andrew",
		"lastName":     "Wansley",
		"emailAddress": "andrew.wansley@gmail.com",
		"age":          int64(32)}
	p := map[string]interface{}{
		"text": "My bio.."}

	ui, _ := tx.Schema().GetInterface("user")
	up, _ := ui.RelationshipByName("profile")
	var createTests = []struct {
		operations []CreateOperation
		output     []CreateCase
	}{
		// Simple Create
		{
			operations: []CreateOperation{
				CreateOperation{
					ModelID: ui.ID(),
					Data:    u,
					Nested:  []NestedOperation{},
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
					ModelID: ui.ID(),
					Data:    u,
					Nested: []NestedOperation{
						NestedCreateOperation{
							Relationship: up,
							Data:         p,
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
					ModelID: up.Target().ID(),
					Data:    p,
					Nested:  []NestedOperation{},
				},
				CreateOperation{
					ModelID: ui.ID(),
					Data:    u,
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
			if v, ok := CreateCase.st["firstName"]; ok {
				_, err := tx.Query(mref, db.Filter(mref, db.Eq("firstName", v))).OneRecord()
				assert.Nil(t, err)
			}
			if v, ok := CreateCase.st["text"]; ok {
				_, err := tx.Query(mref, db.Filter(mref, db.Eq("text", v))).OneRecord()
				assert.Nil(t, err)
			}
		}
	}
}
