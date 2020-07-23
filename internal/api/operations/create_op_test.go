package operations

import (
	"awans.org/aft/internal/db"
	"github.com/go-test/deep"
	"testing"
)

type FindCase struct {
	st        db.Record
	modelName string
}

func TestCreateApply(t *testing.T) {
	appDB := db.NewTest()
	db.AddSampleModels(appDB)
	tx := appDB.NewRWTx()
	u := MakeRecord(tx, "user", `{ 
					"type": "user",
					"firstName":"Andrew",
					"lastName":"Wansley",
					"emailAddress":"andrew.wansley@gmail.com",
					"age": 32}`)
	p := MakeRecord(tx, "profile", `{
		"type":"profile",
		"text": "My bio.."}`)

	var createTests = []struct {
		operations []CreateOperation
		output     []FindCase
	}{
		// Simple Create
		{
			operations: []CreateOperation{
				CreateOperation{
					Record: u,
					Nested: []NestedOperation{},
				},
			},
			output: []FindCase{
				FindCase{
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
							Relationship: db.UserProfile,
							Record:       p,
							Nested:       []NestedOperation{},
						},
					},
				},
			},
			output: []FindCase{
				FindCase{
					st:        u,
					modelName: "user",
				},
				FindCase{
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
							Relationship: db.UserProfile,
							// eventually need this to be a unique prop
							UniqueQuery: UniqueQuery{
								Key: "Text",
								Val: "My bio..",
							},
						},
					},
				},
			},
			output: []FindCase{
				FindCase{
					st:        u,
					modelName: "user",
				},
				FindCase{
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
		for _, findCase := range testCase.output {
			m, _ := tx.Schema().GetModel(findCase.modelName)
			mref := tx.Ref(m.ID())
			found, _ := tx.Query(mref).Filter(mref, db.EqID(findCase.st.ID())).OneRecord()
			if diff := deep.Equal(found, findCase.st); diff != nil {
				t.Error(diff)
			}
		}

	}
}
