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
	AddSampleModels(appDB)
	u := makeStruct(appDB, "user", `{ 
					"type": "user",
					"firstName":"Andrew",
					"lastName":"Wansley", 
					"age": 32}`)
	p := makeStruct(appDB, "profile", `{
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
		appDB = New()
		AddSampleModels(appDB)
		for _, op := range testCase.operations {
			op.Apply(appDB)
		}
		for _, findCase := range testCase.output {
			found, _ := findOneById(appDB, findCase.modelName, getId(findCase.st))
			if diff := deep.Equal(found, findCase.st); diff != nil {
				t.Error(diff)
			}
		}

	}
}
