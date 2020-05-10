package operations

import (
	"awans.org/aft/internal/model"
	"awans.org/aft/internal/server/db"
	"encoding/json"
	"github.com/go-test/deep"
	"github.com/json-iterator/go"
	"testing"
)

func makeRecord(tx db.Tx, modelName string, jsonValue string) model.Record {
	st := tx.MakeRecord(modelName)
	json.Unmarshal([]byte(jsonValue), &st)
	return st
}

func TestParseCreate(t *testing.T) {
	appDB := db.New()
	db.AddSampleModels(appDB)
	tx := appDB.NewRWTx()
	p := Parser{tx: tx}

	var createTests = []struct {
		modelName  string
		jsonString string
		output     interface{}
	}{
		// Simple Create
		{
			modelName: "user",
			jsonString: `{ 
			"firstName":"Andrew",
			"lastName":"Wansley",
			"age": 32}`,
			output: db.CreateOperation{
				Record: makeRecord(tx, "user", `{ 
					"id":"00000000-0000-0000-0000-000000000000",
					"type":"user",
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
			"firstName":"Andrew",
			"lastName":"Wansley",
			"age": 32,
			"profile": {
			  "create": {
			    "text": "My bio.."
			  }
			}}`,
			output: db.CreateOperation{
				Record: makeRecord(tx, "user", `{ 
					"id":"00000000-0000-0000-0000-000000000000",
					"type":"user",
					"firstName":"Andrew",
					"lastName":"Wansley", 
					"age": 32}`),
				Nested: []db.NestedOperation{
					db.NestedCreateOperation{
						Relationship: db.User.Relationships["profile"],
						Record: makeRecord(tx, "profile", `{
							"id":"00000000-0000-0000-0000-000000000000",
							"type":"profile",
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
			"firstName":"Andrew",
			"lastName":"Wansley",
			"age": 32,
			"posts": {
			  "create": [{
			    "text": "post1"
			  }, {
			    "text": "post2"
			  }]
			}}`,
			output: db.CreateOperation{
				Record: makeRecord(tx, "user", `{ 
					"id":"00000000-0000-0000-0000-000000000000",
					"type":"user",
					"firstName":"Andrew",
					"lastName":"Wansley", 
					"age": 32}`),
				Nested: []db.NestedOperation{
					db.NestedCreateOperation{
						Relationship: db.User.Relationships["posts"],
						Record: makeRecord(tx, "post", `{
							"id":"00000000-0000-0000-0000-000000000000",
							"type":"post",
							"text": "post1"}`),
						Nested: []db.NestedOperation{},
					},
					db.NestedCreateOperation{
						Relationship: db.User.Relationships["posts"],
						Record: makeRecord(tx, "post", `{
						    "id":"00000000-0000-0000-0000-000000000000",
						    "type": "post",
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
			"firstName":"Andrew",
			"lastName":"Wansley",
			"age": 32,
			"profile": {
			  "connect": {
			    "id": "57e3f538-d35a-45e8-acdf-0ab916d8194f"
			  }
			}}`,
			output: db.CreateOperation{
				Record: makeRecord(tx, "user", `{ 
					"id":"00000000-0000-0000-0000-000000000000",
					"type": "user",
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
				Record: makeRecord(tx, "user", `{ 
					"id":"00000000-0000-0000-0000-000000000000",
					"type": "user",
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
		parsedOp, err := p.ParseCreate(testCase.modelName, data)
		if err != nil {
			t.Error(err)
		}
		if diff := deep.Equal(parsedOp, testCase.output); diff != nil {
			t.Error(diff)
		}
	}
}
