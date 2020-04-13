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
					"firstName":"Andrew",
					"lastName":"Wansley", 
					"age": 32}`)
	var createTests = []struct {
		input  CreateOperation
		output []FindCase
	}{
		// Simple Create
		{
			input: CreateOperation{
				Struct: u,
				Nested: []NestedOperation{},
			},
			output: []FindCase{
				FindCase{
					st:        u,
					modelName: "user",
				},
			},
		},
	}
	for _, testCase := range createTests {
		testCase.input.Apply(appDB)
		for _, findCase := range testCase.output {
			found := findOneById(appDB, findCase.modelName, getId(u))
			if diff := deep.Equal(found, findCase.st); diff != nil {
				t.Error(diff)
			}
		}

	}
}
