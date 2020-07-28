package parsers

import (
	"awans.org/aft/internal/api/operations"
	"awans.org/aft/internal/db"
	"github.com/google/go-cmp/cmp"
	"github.com/json-iterator/go"
	"testing"
)

func TestParseFindMany(t *testing.T) {
	appDB := db.NewTest()
	db.AddSampleModels(appDB)
	p := Parser{Tx: appDB.NewTx()}

	var findManyTests = []struct {
		modelName  string
		jsonString string
		output     operations.FindManyOperation
	}{
		// Basic String FieldCriterion
		{
			modelName: "user",
			jsonString: `{ 
				"firstName": "Andrew"
			}`,
			output: operations.FindManyOperation{
				ModelID: db.User.ID(),
				FindManyArgs: operations.FindManyArgs{
					Where: operations.Where{
						FieldCriteria: []operations.FieldCriterion{
							operations.FieldCriterion{
								Key: "Firstname",
								Val: "Andrew",
							},
						},
					},
				},
			},
		},
		// Multiple String FieldCriterion
		{
			modelName: "user",
			jsonString: `{ 
				"firstName": "Andrew",
				"lastName": "Wansley",
				"age": 32,
			}`,
			output: operations.FindManyOperation{
				ModelID: db.User.ID(),
				FindManyArgs: operations.FindManyArgs{
					Where: operations.Where{
						FieldCriteria: []operations.FieldCriterion{
							operations.FieldCriterion{
								Key: "Firstname",
								Val: "Andrew",
							},
							operations.FieldCriterion{
								Key: "Lastname",
								Val: "Wansley",
							},
							operations.FieldCriterion{
								Key: "Age",
								Val: int64(32),
							},
						},
					},
				},
			},
		},

		// Single Field To-One Relationship Criterion
		{
			modelName: "user",
			jsonString: `{ 
				"profile": { "text": "This is my bio.." }
			}`,
			output: operations.FindManyOperation{
				ModelID: db.User.ID(),
				FindManyArgs: operations.FindManyArgs{
					Where: operations.Where{
						RelationshipCriteria: []operations.RelationshipCriterion{
							operations.RelationshipCriterion{
								Relationship: db.UserProfile,
								Where: operations.Where{
									FieldCriteria: []operations.FieldCriterion{
										operations.FieldCriterion{
											Key: "Text",
											Val: "This is my bio..",
										},
									},
								},
							},
						},
					},
				},
			},
		},

		// Single Field To-One Relationship Criterion
		// with Nested Relationship Criterion
		{
			modelName: "user",
			jsonString: `{
						"profile": {
							"text": "This is my bio..",
							"user": {
							  "firstName": "Andrew"
							}
						}
					}`,
			output: operations.FindManyOperation{
				ModelID: db.User.ID(),
				FindManyArgs: operations.FindManyArgs{
					Where: operations.Where{
						RelationshipCriteria: []operations.RelationshipCriterion{
							operations.RelationshipCriterion{
								Relationship: db.UserProfile,
								Where: operations.Where{
									RelationshipCriteria: []operations.RelationshipCriterion{
										operations.RelationshipCriterion{
											Relationship: db.ProfileUser,
											Where: operations.Where{
												FieldCriteria: []operations.FieldCriterion{
													operations.FieldCriterion{
														Key: "Firstname",
														Val: "Andrew",
													},
												},
											},
										},
									},
									FieldCriteria: []operations.FieldCriterion{
										operations.FieldCriterion{
											Key: "Text",
											Val: "This is my bio..",
										},
									},
								},
							},
						},
					},
				},
			},
		},

		// Single Field To-Many "Some" Relationship Criterion
		{
			modelName:  "user",
			jsonString: `{ "posts": { "some": { "text": "This is my bio.." } } }`,
			output: operations.FindManyOperation{
				ModelID: db.User.ID(),
				FindManyArgs: operations.FindManyArgs{
					Where: operations.Where{
						AggregateRelationshipCriteria: []operations.AggregateRelationshipCriterion{
							operations.AggregateRelationshipCriterion{
								Aggregation: db.Some,
								RelationshipCriterion: operations.RelationshipCriterion{
									Relationship: db.UserPosts,
									Where: operations.Where{
										FieldCriteria: []operations.FieldCriterion{
											operations.FieldCriterion{
												Key: "Text",
												Val: "This is my bio..",
											},
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
	for _, testCase := range findManyTests {
		data := make(map[string]interface{})
		var where map[string]interface{}
		jsoniter.Unmarshal([]byte(testCase.jsonString), &where)
		data["where"] = where
		parsedOp, err := p.ParseFindMany(testCase.modelName, data)
		if err != nil {
			t.Error(err)
		}
		opts := append(CmpOpts(), IgnoreRecIDs)
		diff := cmp.Diff(testCase.output, parsedOp, opts...)
		if diff != "" {
			t.Errorf("(-want +got):\n%s", diff)
		}
	}
}
