package parsers

import (
	"awans.org/aft/internal/api/operations"
	"awans.org/aft/internal/db"
	"github.com/google/go-cmp/cmp"
	"github.com/json-iterator/go"
	"testing"
)

func TestParseDelete(t *testing.T) {
	appDB := db.NewTest()
	db.AddSampleModels(appDB)
	p := Parser{Tx: appDB.NewTx()}

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
				Where:   operations.Where{},
				ModelID: db.Profile.ID(),
				Include: operations.Include{},
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
				Where:   operations.Where{},
				ModelID: db.User.ID(),
				Include: operations.Include{},
				Nested: []operations.NestedOperation{
					operations.NestedDeleteOperation{
						Relationship: db.UserProfile,
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
				Where:   operations.Where{},
				ModelID: db.Profile.ID(),
				Include: operations.Include{},
				Nested: []operations.NestedOperation{
					operations.NestedDeleteOperation{
						Relationship: db.ProfileUser,
						Nested: []operations.NestedOperation{
							operations.NestedDeleteManyOperation{
								Relationship: db.UserPosts,
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
				Where:   operations.Where{},
				ModelID: db.Profile.ID(),
				Include: operations.Include{},
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
						Relationship: db.ProfileUser,
						Nested: []operations.NestedOperation{
							operations.NestedDeleteManyOperation{
								Relationship: db.UserPosts,
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
				Where:   operations.Where{},
				ModelID: db.Profile.ID(),
				Include: operations.Include{},
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
						Relationship: db.ProfileUser,
						Nested: []operations.NestedOperation{
							operations.NestedDeleteManyOperation{
								Relationship: db.UserPosts,
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
