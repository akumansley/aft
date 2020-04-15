package db

import (
	"encoding/json"
	"github.com/go-test/deep"
	"testing"
)

func makeStruct(db DB, modelName string, jsonValue string) interface{} {
	st := db.MakeStruct(modelName)
	json.Unmarshal([]byte(jsonValue), &st)
	return st
}

type FindCase struct {
	st        interface{}
	modelName string
}

func TestCreateApply(t *testing.T) {
	appDB := New()
	appDB.AddSampleModels()
	u := makeStruct(appDB, "user", `{ 
					"id":"15852d31-3bd4-4fc4-abd0-e4c7497644ab",
					"type": "user",
					"firstName":"Andrew",
					"lastName":"Wansley", 
					"age": 32}`)
	p := makeStruct(appDB, "profile", `{
		"id": "c8f857ca-204c-46ab-a96e-d69c1df2fa4f",
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
					Struct: u,
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
					Struct: u,
					Nested: []NestedOperation{
						NestedCreateOperation{
							Relationship: User.Relationships["profile"],
							Struct:       p,
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
					Struct: p,
					Nested: []NestedOperation{},
				},
				CreateOperation{
					Struct: u,
					Nested: []NestedOperation{
						NestedConnectOperation{
							Relationship: User.Relationships["profile"],
							UniqueQuery: UniqueQuery{
								Key: "Id",
								Val: getId(p),
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
		appDB = New()
		appDB.AddSampleModels()
		for _, op := range testCase.operations {
			op.Apply(appDB)
		}
		for _, findCase := range testCase.output {
			found := findOneById(appDB, findCase.modelName, getId(findCase.st))
			if diff := deep.Equal(found, findCase.st); diff != nil {
				t.Error(diff)
			}
		}

	}
}
