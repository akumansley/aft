package parsers

import (
	"testing"

	"awans.org/aft/internal/api/operations"
	"awans.org/aft/internal/db"
	"github.com/google/go-cmp/cmp"
	"github.com/google/uuid"
	jsoniter "github.com/json-iterator/go"
)

func TestParseCreate(t *testing.T) {
	appDB := db.NewTest()
	db.AddSampleModels(appDB)
	tx := appDB.NewRWTx()
	p := Parser{Tx: tx}

	var createTests = []struct {
		modelName  string
		jsonString string
		output     interface{}
	}{
		// Simple Create
		{
			modelName: "user",
			jsonString: `{"data": { 
			"firstName":"Andrew",
			"lastName":"Wansley",
			"age": 32,
			"emailAddress":"andrew.wansley@gmail.com"}}`,
			output: operations.CreateOperation{
				ModelID: db.User.ID(),
				Data: map[string]interface{}{
					"firstName":    "Andrew",
					"lastName":     "Wansley",
					"emailAddress": "andrew.wansley@gmail.com",
					"age":          32.0},
				Nested: []operations.NestedOperation{},
			},
		},
		// Nested Single Create
		{
			modelName: "user",
			jsonString: `{"data" : {
			"firstName":"Andrew",
			"lastName":"Wansley",
			"age": 32,
			"emailAddress":"andrew.wansley@gmail.com",
			"profile": {
			  "create": {
			    "text": "My bio.."
			  }
			}}}`,
			output: operations.CreateOperation{
				ModelID: db.User.ID(),
				Data: map[string]interface{}{
					"firstName":    "Andrew",
					"lastName":     "Wansley",
					"emailAddress": "andrew.wansley@gmail.com",
					"age":          32.0},
				Nested: []operations.NestedOperation{
					operations.NestedCreateOperation{
						Model:        db.Profile,
						Relationship: db.UserProfile,
						Data: map[string]interface{}{
							"text": "My bio.."},
						Nested: []operations.NestedOperation{},
					},
				},
			},
		},
		// Nested Multiple Create
		{
			modelName: "user",
			jsonString: `{"data" :{
			"firstName":"Andrew",
			"lastName":"Wansley",
			"age": 32,
			"emailAddress":"andrew.wansley@gmail.com",
			"posts": {
			  "create": [{
			    "text": "post1"
			  }, {
			    "text": "post2"
			  }]
			}}}`,
			output: operations.CreateOperation{
				ModelID: db.User.ID(),
				Data: map[string]interface{}{
					"firstName":    "Andrew",
					"lastName":     "Wansley",
					"emailAddress": "andrew.wansley@gmail.com",
					"age":          32.0},
				Nested: []operations.NestedOperation{
					operations.NestedCreateOperation{
						Relationship: db.UserPosts,
						Model:        db.Post,
						Data: map[string]interface{}{
							"text": "post1"},
						Nested: []operations.NestedOperation{},
					},
					operations.NestedCreateOperation{
						Relationship: db.UserPosts,
						Model:        db.Post,
						Data: map[string]interface{}{
							"text": "post2"},
						Nested: []operations.NestedOperation{},
					},
				},
			},
		},
		// Nested Connect
		{
			modelName: "user",
			jsonString: `{"data" : {
			"firstName":"Andrew",
			"lastName":"Wansley",
			"age": 32,
			"emailAddress":"andrew.wansley@gmail.com",
			"profile": {
			  "connect": {
			    "id": "57e3f538-d35a-45e8-acdf-0ab916d8194f"
			  }
			}}}`,
			output: operations.CreateOperation{
				ModelID: db.User.ID(),
				Data: map[string]interface{}{
					"firstName":    "Andrew",
					"lastName":     "Wansley",
					"emailAddress": "andrew.wansley@gmail.com",
					"age":          32.0},
				Nested: []operations.NestedOperation{
					operations.NestedConnectOperation{
						Relationship: db.UserProfile,
						Where: operations.Where{
							FieldCriteria: []operations.FieldCriterion{
								operations.FieldCriterion{
									Key: "Id",
									Val: uuid.MustParse("57e3f538-d35a-45e8-acdf-0ab916d8194f"),
								},
							},
						},
					},
				},
			},
		},
		// Nested Multi Connect
		{
			modelName: "user",
			jsonString: `{"data" : {
			"firstName":"Andrew",
			"lastName":"Wansley",
			"age": 32,
			"emailAddress":"andrew.wansley@gmail.com",
			"posts": {
			  "connect": [{
			    "id": "57e3f538-d35a-45e8-acdf-0ab916d8194f"
			  }, {
			    "id": "6327fe0e-c936-4332-85cd-f1b42f6f337a",
			  }]
			}}}`,
			output: operations.CreateOperation{
				ModelID: db.User.ID(),
				Data: map[string]interface{}{
					"firstName":    "Andrew",
					"lastName":     "Wansley",
					"emailAddress": "andrew.wansley@gmail.com",
					"age":          32.0},
				Nested: []operations.NestedOperation{
					operations.NestedConnectOperation{
						Relationship: db.UserPosts,
						Where: operations.Where{
							FieldCriteria: []operations.FieldCriterion{
								operations.FieldCriterion{
									Key: "Id",
									Val: uuid.MustParse("57e3f538-d35a-45e8-acdf-0ab916d8194f"),
								},
							},
						},
					},
					operations.NestedConnectOperation{
						Relationship: db.UserPosts,
						Where: operations.Where{
							FieldCriteria: []operations.FieldCriterion{
								operations.FieldCriterion{
									Key: "Id",
									Val: uuid.MustParse("6327fe0e-c936-4332-85cd-f1b42f6f337a"),
								},
							},
						},
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
			t.Fatal(err)
		}

		opts := append(CmpOpts(), IgnoreRecIDs)

		diff := cmp.Diff(testCase.output, parsedOp, opts...)
		if diff != "" {
			t.Errorf("(-want +got):\n%s", diff)
		}
	}
}
