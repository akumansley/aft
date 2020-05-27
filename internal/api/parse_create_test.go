package api

import (
	"awans.org/aft/internal/db"
	"github.com/go-test/deep"
	"github.com/google/uuid"
	"github.com/json-iterator/go"
	"testing"
)

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
			"age": 32,
			"email":"andrew.wansley@gmail.com"}`,
			output: CreateOperation{
				Record: makeRecord(tx, "user", `{ 
					"id":"00000000-0000-0000-0000-000000000000",
					"firstName":"Andrew",
					"lastName":"Wansley",
					"email":"andrew.wansley@gmail.com",
					"age": 32}`),
				Nested: []NestedOperation{},
			},
		},
		// Nested Single Create
		{
			modelName: "user",
			jsonString: `{
			"firstName":"Andrew",
			"lastName":"Wansley",
			"age": 32,
			"email":"andrew.wansley@gmail.com",
			"profile": {
			  "create": {
			    "text": "My bio.."
			  }
			}}`,
			output: CreateOperation{
				Record: makeRecord(tx, "user", `{ 
					"id":"00000000-0000-0000-0000-000000000000",
					"firstName":"Andrew",
					"lastName":"Wansley",
					"email":"andrew.wansley@gmail.com",
					"age": 32}`),
				Nested: []NestedOperation{
					NestedCreateOperation{
						Binding: db.UserProfile.Left(),
						Record: makeRecord(tx, "profile", `{
							"id":"00000000-0000-0000-0000-000000000000",
							"text": "My bio.."}`),
						Nested: []NestedOperation{},
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
			"email":"andrew.wansley@gmail.com",
			"posts": {
			  "create": [{
			    "text": "post1"
			  }, {
			    "text": "post2"
			  }]
			}}`,
			output: CreateOperation{
				Record: makeRecord(tx, "user", `{ 
					"id":"00000000-0000-0000-0000-000000000000",
					"firstName":"Andrew",
					"lastName":"Wansley",
					"email":"andrew.wansley@gmail.com",
					"age": 32}`),
				Nested: []NestedOperation{
					NestedCreateOperation{
						Binding: db.UserPosts.Left(),
						Record: makeRecord(tx, "post", `{
							"id":"00000000-0000-0000-0000-000000000000",
							"text": "post1"}`),
						Nested: []NestedOperation{},
					},
					NestedCreateOperation{
						Binding: db.UserPosts.Left(),
						Record: makeRecord(tx, "post", `{
						    "id":"00000000-0000-0000-0000-000000000000",
						    "text": "post2"}`),
						Nested: []NestedOperation{},
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
			"email":"andrew.wansley@gmail.com",
			"profile": {
			  "connect": {
			    "id": "57e3f538-d35a-45e8-acdf-0ab916d8194f"
			  }
			}}`,
			output: CreateOperation{
				Record: makeRecord(tx, "user", `{ 
					"id":"00000000-0000-0000-0000-000000000000",
					"firstName":"Andrew",
					"lastName":"Wansley", 
					"email":"andrew.wansley@gmail.com",
					"age": 32}`),
				Nested: []NestedOperation{
					NestedConnectOperation{
						Binding: db.UserProfile.Left(),
						UniqueQuery: UniqueQuery{
							Key: "id",
							Val: uuid.MustParse("57e3f538-d35a-45e8-acdf-0ab916d8194f")},
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
			"email":"andrew.wansley@gmail.com",
			"posts": {
			  "connect": [{
			    "id": "57e3f538-d35a-45e8-acdf-0ab916d8194f"
			  }, {
			    "id": "6327fe0e-c936-4332-85cd-f1b42f6f337a",
			  }]
			}}`,
			output: CreateOperation{
				Record: makeRecord(tx, "user", `{ 
					"id":"00000000-0000-0000-0000-000000000000",
					"firstName":"Andrew",
					"lastName":"Wansley", 
					"email":"andrew.wansley@gmail.com",
					"age": 32}`),
				Nested: []NestedOperation{
					NestedConnectOperation{
						Binding: db.UserPosts.Left(),
						UniqueQuery: UniqueQuery{
							Key: "id",
							Val: uuid.MustParse("57e3f538-d35a-45e8-acdf-0ab916d8194f")},
					},
					NestedConnectOperation{
						Binding: db.UserPosts.Left(),
						UniqueQuery: UniqueQuery{
							Key: "id",
							Val: uuid.MustParse("6327fe0e-c936-4332-85cd-f1b42f6f337a")},
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
