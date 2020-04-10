package operations

import (
	"awans.org/aft/internal/server/db"
	"encoding/json"
	"github.com/go-test/deep"
	"github.com/json-iterator/go"
	"testing"
)

func makeStruct(db db.DB, modelName string, jsonValue string) interface{} {
	st := db.MakeStruct(modelName)
	json.Unmarshal([]byte(jsonValue), &st)
	return st
}

func TestParseCreate(t *testing.T) {
	appDB := db.New()
	appDB.AddSampleModels()
	p := Parser{db: appDB}

	var createTests = []struct {
		modelName  string
		jsonString string
		output     interface{}
	}{
		// Simple Create
		{
			modelName: "user",
			jsonString: `{ 
			"id":"15852d31-3bd4-4fc4-abd0-e4c7497644ab",
			"firstName":"Andrew",
			"lastName":"Wansley",
			"age": 32}`,
			output: db.CreateOperation{
				Struct: makeStruct(appDB, "user", `{ 
					"id":"15852d31-3bd4-4fc4-abd0-e4c7497644ab",
					"firstName":"Andrew",
					"lastName":"Wansley", 
					"age": 32}`),
				Nested: []db.NestedOperation{},
			},
		},
		// Nested Single Create
		{
			modelName: "user",
			jsonString: `{
			"id":"15852d31-3bd4-4fc4-abd0-e4c7497644ab",
			"firstName":"Andrew",
			"lastName":"Wansley",
			"age": 32,
			"profile": {
			  "create": {
			    "id": "c8f857ca-204c-46ab-a96e-d69c1df2fa4f",
			    "text": "My bio.."
			  }
			}}`,
			output: db.CreateOperation{
				Struct: makeStruct(appDB, "user", `{ 
					"id":"15852d31-3bd4-4fc4-abd0-e4c7497644ab",
					"firstName":"Andrew",
					"lastName":"Wansley", 
					"age": 32}`),
				Nested: []db.NestedOperation{
					db.NestedCreateOperation{
						Relationship: db.User.Relationships["profile"],
						Struct: makeStruct(appDB, "profile", `{
						    "id": "c8f857ca-204c-46ab-a96e-d69c1df2fa4f",
						    "text": "My bio.."}`),
						Nested: []db.NestedOperation{},
					},
				},
			},
		},
		// Nested Multiple Create
		{
			modelName: "user",
			jsonString: `{
			"id":"15852d31-3bd4-4fc4-abd0-e4c7497644ab",
			"firstName":"Andrew",
			"lastName":"Wansley",
			"age": 32,
			"posts": {
			  "create": [{
			    "id": "57e3f538-d35a-45e8-acdf-0ab916d8194f",
			    "text": "post1"
			  }, {
			    "id": "6327fe0e-c936-4332-85cd-f1b42f6f337a",
			    "text": "post2"
			  }]
			}}`,
			output: db.CreateOperation{
				Struct: makeStruct(appDB, "user", `{ 
					"id":"15852d31-3bd4-4fc4-abd0-e4c7497644ab",
					"firstName":"Andrew",
					"lastName":"Wansley", 
					"age": 32}`),
				Nested: []db.NestedOperation{
					db.NestedCreateOperation{
						Relationship: db.User.Relationships["posts"],
						Struct: makeStruct(appDB, "post", `{
						    "id": "57e3f538-d35a-45e8-acdf-0ab916d8194f",
						    "text": "post1"}`),
						Nested: []db.NestedOperation{},
					},
					db.NestedCreateOperation{
						Relationship: db.User.Relationships["posts"],
						Struct: makeStruct(appDB, "post", `{
						    "id": "6327fe0e-c936-4332-85cd-f1b42f6f337a",
						    "text": "post2"}`),
						Nested: []db.NestedOperation{},
					},
				},
			},
		},
		// Nested Connect
		{
			modelName: "user",
			jsonString: `{
			"id":"15852d31-3bd4-4fc4-abd0-e4c7497644ab",
			"firstName":"Andrew",
			"lastName":"Wansley",
			"age": 32,
			"profile": {
			  "connect": {
			    "id": "57e3f538-d35a-45e8-acdf-0ab916d8194f"
			  }
			}}`,
			output: db.CreateOperation{
				Struct: makeStruct(appDB, "user", `{ 
					"id":"15852d31-3bd4-4fc4-abd0-e4c7497644ab",
					"firstName":"Andrew",
					"lastName":"Wansley", 
					"age": 32}`),
				Nested: []db.NestedOperation{
					db.NestedConnectOperation{
						Relationship: db.User.Relationships["profile"],
						UniqueQuery: db.UniqueQuery{
							Key: "id",
							Val: "57e3f538-d35a-45e8-acdf-0ab916d8194f"},
					},
				},
			},
		},
		// Nested Multi Connect
		{
			modelName: "user",
			jsonString: `{
			"id":"15852d31-3bd4-4fc4-abd0-e4c7497644ab",
			"firstName":"Andrew",
			"lastName":"Wansley",
			"age": 32,
			"posts": {
			  "connect": [{
			    "id": "57e3f538-d35a-45e8-acdf-0ab916d8194f"
			  }, {
			    "id": "6327fe0e-c936-4332-85cd-f1b42f6f337a",
			  }]
			}}`,
			output: db.CreateOperation{
				Struct: makeStruct(appDB, "user", `{ 
					"id":"15852d31-3bd4-4fc4-abd0-e4c7497644ab",
					"firstName":"Andrew",
					"lastName":"Wansley", 
					"age": 32}`),
				Nested: []db.NestedOperation{
					db.NestedConnectOperation{
						Relationship: db.User.Relationships["posts"],
						UniqueQuery: db.UniqueQuery{
							Key: "id",
							Val: "57e3f538-d35a-45e8-acdf-0ab916d8194f"},
					},
					db.NestedConnectOperation{
						Relationship: db.User.Relationships["posts"],
						UniqueQuery: db.UniqueQuery{
							Key: "id",
							Val: "6327fe0e-c936-4332-85cd-f1b42f6f337a"},
					},
				},
			},
		},
	}
	for _, testCase := range createTests {
		var data map[string]interface{}
		jsoniter.Unmarshal([]byte(testCase.jsonString), &data)
		parsedOp := p.ParseCreate(testCase.modelName, data)
		if diff := deep.Equal(parsedOp, testCase.output); diff != nil {
			t.Error(diff)
		}
	}
}
