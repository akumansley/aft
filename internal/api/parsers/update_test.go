package parsers

import (
	"awans.org/aft/internal/api/operations"
	"awans.org/aft/internal/db"
	"github.com/google/go-cmp/cmp"
	"github.com/json-iterator/go"
	"testing"
)

func TestParseUpdate(t *testing.T) {
	appDB := db.NewTest()
	db.AddSampleModels(appDB)
	tx := appDB.NewRWTx()
	p := Parser{Tx: tx}

	var updateTests = []struct {
		modelName  string
		jsonString string
		output     interface{}
	}{
		// Simple Update
		{
			modelName: "user",
			jsonString: `{
				"data": {"firstName":"Chase"}, 
				"where" : {"firstName" : "Andrew"}
			}`,
			output: operations.UpdateOperation{
				ModelID: db.User.ID(),
				FindArgs: operations.FindArgs{
					Where: operations.Where{
						FieldCriteria: []operations.FieldCriterion{
							operations.FieldCriterion{
								Key: "Firstname",
								Val: "Andrew",
							},
						},
					},
					Include: operations.Include{},
				},
				Data:   map[string]interface{}{"firstName": "Chase"},
				Nested: []operations.NestedOperation{},
			},
		},

		// Nested Delete
		{
			modelName: "user",
			jsonString: `{
				"data": {"firstName":"Chase", "profile" : {"delete" : true}}, 
				"where" : {"firstName" : "Andrew"}
			}`,
			output: operations.UpdateOperation{
				ModelID: db.User.ID(),
				FindArgs: operations.FindArgs{
					Where: operations.Where{
						FieldCriteria: []operations.FieldCriterion{
							operations.FieldCriterion{
								Key: "Firstname",
								Val: "Andrew",
							},
						},
					},
					Include: operations.Include{},
				},
				Data: map[string]interface{}{"firstName": "Chase"},
				Nested: []operations.NestedOperation{
					operations.NestedDeleteOperation{
						Relationship: db.UserProfile,
					},
				},
			},
		},

		// Nested DeleteMany
		{
			modelName: "user",
			jsonString: `{
				"data": { "firstName":"Chase", "profile" : {"deleteMany" : true}}, 
				"where" : {"firstName" : "Andrew"}
			}`,
			output: operations.UpdateOperation{
				ModelID: db.User.ID(),
				FindArgs: operations.FindArgs{
					Where: operations.Where{
						FieldCriteria: []operations.FieldCriterion{
							operations.FieldCriterion{
								Key: "Firstname",
								Val: "Andrew",
							},
						},
					},
					Include: operations.Include{},
				},
				Data: map[string]interface{}{"firstName": "Chase"},
				Nested: []operations.NestedOperation{
					operations.NestedDeleteManyOperation{
						Where:        operations.Where{},
						Relationship: db.UserProfile,
					},
				},
			},
		},

		// Nested UpdateMany
		{
			modelName: "user",
			jsonString: `{
				"data": { 
					"firstName":"Chase", 
					"profile" : {
						"updateMany" : {
							"data" : {"text" : "cool"}
						}
					}
				},
				"where" : {"firstName" : "Andrew"}}`,
			output: operations.UpdateOperation{
				ModelID: db.User.ID(),
				FindArgs: operations.FindArgs{
					Where: operations.Where{
						FieldCriteria: []operations.FieldCriterion{
							operations.FieldCriterion{
								Key: "Firstname",
								Val: "Andrew",
							},
						},
					},
					Include: operations.Include{},
				},
				Data: map[string]interface{}{"firstName": "Chase"},
				Nested: []operations.NestedOperation{
					operations.NestedUpdateManyOperation{
						Relationship: db.UserProfile,
						Data:         map[string]interface{}{"text": "cool"},
						Nested:       []operations.NestedOperation{},
					},
				},
			},
		},
	}
	for _, testCase := range updateTests {
		var data map[string]interface{}
		jsoniter.Unmarshal([]byte(testCase.jsonString), &data)
		parsedOp, err := p.ParseUpdate(testCase.modelName, data)
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
