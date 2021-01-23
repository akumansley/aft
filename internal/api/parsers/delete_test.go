package parsers

import (
	"testing"

	"awans.org/aft/internal/api/operations"
	"awans.org/aft/internal/db"
	"github.com/google/go-cmp/cmp"
	jsoniter "github.com/json-iterator/go"
)

func TestParseDelete(t *testing.T) {
	appDB := db.NewTest()
	db.AddSampleModels(appDB)
	tx := appDB.NewTx()
	p := Parser{Tx: tx}

	var deletionTests = []struct {
		modelName  string
		jsonString string
		output     operations.DeleteOperation
	}{
		// Simple delete
		{
			modelName:  "profile",
			jsonString: `{}`,
			output: operations.DeleteOperation{
				FindArgs: operations.FindArgs{
					Where:   operations.Where{},
					Include: operations.Include{},
				},
				ModelID: db.Profile.ID(),
				Nested:  nil,
			},
		},

		// Nested delete
		{
			modelName: "user",
			jsonString: `{
				"delete" : {"profile" : true}
			}`,
			output: operations.DeleteOperation{
				ModelID: db.User.ID(),
				FindArgs: operations.FindArgs{
					Where:   operations.Where{},
					Include: operations.Include{},
				},
				Nested: []operations.NestedOperation{
					operations.NestedDeleteOperation{
						Relationship: db.UserProfile.Load(tx),
					},
				},
			},
		},
		// Double Nested delete
		{
			modelName: "profile",
			jsonString: `{ 
			   "delete" : {
			   		"user" : {
			   			"deleteMany" : {"posts" : true}
			   		}
			   	}
			}`,
			output: operations.DeleteOperation{
				ModelID: db.Profile.ID(),
				FindArgs: operations.FindArgs{
					Where:   operations.Where{},
					Include: operations.Include{},
				},
				Nested: []operations.NestedOperation{
					operations.NestedDeleteOperation{
						Relationship: db.ProfileUser.Load(tx),
						Nested: []operations.NestedOperation{
							operations.NestedDeleteManyOperation{
								Relationship: db.UserPosts.Load(tx),
							},
						},
					},
				},
			},
		},

		// Double Nested delete with a where
		{
			modelName: "profile",
			jsonString: `{ 
			   "delete" : {
			   		"user" : {
			   			"where" : {"firstName" : "john"}, 
			   			"deleteMany" : {"posts" : true}
			   		}
			   	}
			}`,
			output: operations.DeleteOperation{
				ModelID: db.Profile.ID(),
				FindArgs: operations.FindArgs{
					Where:   operations.Where{},
					Include: operations.Include{},
				},
				Nested: []operations.NestedOperation{
					operations.NestedDeleteOperation{
						Where: operations.Where{
							FieldCriteria: []operations.FieldCriterion{
								operations.FieldCriterion{
									Key: "Firstname",
									Val: "john",
								},
							},
						},
						Relationship: db.ProfileUser.Load(tx),
						Nested: []operations.NestedOperation{
							operations.NestedDeleteManyOperation{
								Relationship: db.UserPosts.Load(tx),
							},
						},
					},
				},
			},
		},

		// Double Nested delete with a nested where
		{
			modelName: "profile",
			jsonString: `{ 
			   "delete" : {
			   		"user" : {
			   			"where" : {"firstName" : "john"}, 
			   			"deleteMany" : {
			   				"posts" : {
			   					"where" : {"text" : ""}
			   				}
			   			}
			   		}
			   	}
			}`,
			output: operations.DeleteOperation{
				ModelID: db.Profile.ID(),
				FindArgs: operations.FindArgs{
					Where:   operations.Where{},
					Include: operations.Include{},
				},
				Nested: []operations.NestedOperation{
					operations.NestedDeleteOperation{
						Where: operations.Where{
							FieldCriteria: []operations.FieldCriterion{
								operations.FieldCriterion{
									Key: "Firstname",
									Val: "john",
								},
							},
						},
						Relationship: db.ProfileUser.Load(tx),
						Nested: []operations.NestedOperation{
							operations.NestedDeleteManyOperation{
								Relationship: db.UserPosts.Load(tx),
								Where: operations.Where{
									FieldCriteria: []operations.FieldCriterion{
										operations.FieldCriterion{
											Key: "Text",
											Val: "",
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}
	for _, testCase := range deletionTests {
		var data map[string]interface{}
		jsoniter.Unmarshal([]byte(testCase.jsonString), &data)
		parsedOp, err := p.ParseDelete(testCase.modelName, data)
		if err != nil {
			t.Error(err)
		}
		diff := cmp.Diff(testCase.output, parsedOp, CmpOpts()...)
		if diff != "" {
			t.Errorf("(-want +got):\n%s", diff)
		}
	}
}
