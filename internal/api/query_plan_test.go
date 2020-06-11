package api

// import (
// 	"awans.org/aft/internal/db"
// 	"fmt"
// 	"github.com/go-test/deep"
// 	"testing"
// )

// func TestQueryPlanner(t *testing.T) {
// 	userFC := FieldCriterion{
// 		Key: "Firstname",
// 		Val: "Andrew",
// 	}

// 	profileFC := FieldCriterion{
// 		Key: "Text",
// 		Val: "This is my bio..",
// 	}
// 	var qpTests = []struct {
// 		operation FindManyOperation
// 		output    relation
// 	}{
// 		{
// 			operation: FindManyOperation{
// 				ModelID: db.User.ID,
// 				Where: Where{
// 					FieldCriteria: []FieldCriterion{
// 						userFC,
// 					},
// 				},
// 			},
// 			output: &seqscan{db.User.ID, []db.Matcher{userFC.Matcher()}},
// 		},
// 		{
// 			operation: FindManyOperation{
// 				ModelID: db.User.ID,
// 				Where: Where{
// 					FieldCriteria: []FieldCriterion{
// 						userFC,
// 					},
// 					RelationshipCriteria: []RelationshipCriterion{
// 						RelationshipCriterion{
// 							Binding: db.UserProfile.Left(),
// 							Where: Where{
// 								FieldCriteria: []FieldCriterion{
// 									profileFC,
// 								},
// 							},
// 						},
// 					},
// 				},
// 			},
// 			output: &join{
// 				&seqscan{db.User.ID, []db.Matcher{userFC.Matcher()}},
// 				&seqscan{db.Profile.ID, []db.Matcher{profileFC.Matcher()}},
// 				db.UserProfile.Left(),
// 			},
// 		},
// 	}

// 	for _, testCase := range qpTests {
// 		r := Plan(testCase.operation)
// 		if diff := deep.Equal(r, testCase.output); diff != nil {
// 			fmt.Printf("%v\n", r)
// 			t.Error(diff)
// 		}
// 	}

// }
