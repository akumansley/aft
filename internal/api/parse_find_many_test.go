package api

import (
	"awans.org/aft/internal/db"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/json-iterator/go"
	"testing"
)

func TestParseFindMany(t *testing.T) {
	appDB := db.NewTest()
	db.AddSampleModels(appDB)
	p := Parser{tx: appDB.NewTx()}

	var findManyTests = []struct {
		modelName  string
		jsonString string
		output     FindManyOperation
	}{
		// Basic String FieldCriterion
		{
			modelName: "user",
			jsonString: `{ 
				"firstName": "Andrew"
			}`,
			output: FindManyOperation{
				ModelID: db.User.ID,
				Where: Where{
					FieldCriteria: []FieldCriterion{
						FieldCriterion{
							Key: "Firstname",
							Val: "Andrew",
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
			output: FindManyOperation{
				ModelID: db.User.ID,
				Where: Where{
					FieldCriteria: []FieldCriterion{
						FieldCriterion{
							Key: "Firstname",
							Val: "Andrew",
						},
						FieldCriterion{
							Key: "Lastname",
							Val: "Wansley",
						},
						FieldCriterion{
							Key: "Age",
							Val: int64(32),
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
			output: FindManyOperation{
				ModelID: db.User.ID,
				Where: Where{
					RelationshipCriteria: []RelationshipCriterion{
						RelationshipCriterion{
							Binding: db.UserProfile.Left(),
							Where: Where{
								FieldCriteria: []FieldCriterion{
									FieldCriterion{
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
			output: FindManyOperation{
				ModelID: db.User.ID,
				Where: Where{
					RelationshipCriteria: []RelationshipCriterion{
						RelationshipCriterion{
							Binding: db.UserProfile.Left(),
							Where: Where{
								RelationshipCriteria: []RelationshipCriterion{
									RelationshipCriterion{
										Binding: db.UserProfile.Right(),
										Where: Where{
											FieldCriteria: []FieldCriterion{
												FieldCriterion{
													Key: "Firstname",
													Val: "Andrew",
												},
											},
										},
									},
								},
								FieldCriteria: []FieldCriterion{
									FieldCriterion{
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

		// Single Field To-Many "Some" Relationship Criterion
		{
			modelName:  "user",
			jsonString: `{ "posts": { "some": { "text": "This is my bio.." } } }`,
			output: FindManyOperation{
				ModelID: db.User.ID,
				Where: Where{
					AggregateRelationshipCriteria: []AggregateRelationshipCriterion{
						AggregateRelationshipCriterion{
							Aggregation: Some,
							RelationshipCriterion: RelationshipCriterion{
								Binding: db.UserPosts.Left(),
								Where: Where{
									FieldCriteria: []FieldCriterion{
										FieldCriterion{
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
	}
	for _, testCase := range findManyTests {
		var data map[string]interface{}
		jsoniter.Unmarshal([]byte(testCase.jsonString), &data)
		parsedOp, err := p.ParseFindMany(testCase.modelName, data)
		if err != nil {
			t.Error(err)
		}
		tFC := cmpopts.SortSlices(func(a, b FieldCriterion) bool {
			return a.Key < b.Key
		})
		tRC := cmpopts.SortSlices(func(a, b RelationshipCriterion) bool {
			return a.Binding.Name() < b.Binding.Name()
		})
		diff := cmp.Diff(testCase.output, parsedOp, tFC, tRC)
		if diff != "" {
			t.Errorf("(-want +got):\n%s", diff)
		}
	}
}
