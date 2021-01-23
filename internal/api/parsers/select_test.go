package parsers

import (
	"testing"

	"awans.org/aft/internal/api"
	"awans.org/aft/internal/api/operations"
	"awans.org/aft/internal/db"
	"github.com/google/go-cmp/cmp"
	jsoniter "github.com/json-iterator/go"
)

func TestParseSelect(t *testing.T) {
	appDB := db.NewTest()
	db.AddSampleModels(appDB)
	tx := appDB.NewTx()
	p := Parser{Tx: tx}

	up, _ := tx.Schema().GetRelationshipByID(db.UserProfile.ID())
	userModel, _ := tx.Schema().GetModelByID(db.User.ID())
	var inclusionTests = []struct {
		model      db.Interface
		jsonString string
		output     operations.Select
	}{
		// Simple Select
		{
			model: userModel,
			jsonString: `{
			   "firstName" : true,
			   "lastName" : true,
			   "profile": true
			}`,
			output: operations.Select{
				Selecting: true,
				Fields:    api.Set{"firstName": api.Void{}, "lastName": api.Void{}},
				Selects: []operations.Selection{
					operations.Selection{
						Relationship: up,
					},
				},
			},
		},

		// Simple Select with where
		{
			model: userModel,
			jsonString: `{
			   "firstName" : true,
			   "profile": {"where" : {"text" : "mybio..."}}
			}`,
			output: operations.Select{
				Selecting: true,
				Fields:    api.Set{"firstName": api.Void{}},
				Selects: []operations.Selection{
					operations.Selection{
						Relationship: up,
						NestedFindMany: operations.FindArgs{
							Where: operations.Where{
								FieldCriteria: []operations.FieldCriterion{
									operations.FieldCriterion{
										Key: "Text",
										Val: "mybio...",
									},
								},
							},
						},
					},
				},
			},
		},
	}
	for _, testCase := range inclusionTests {
		var data map[string]interface{}
		jsoniter.Unmarshal([]byte(testCase.jsonString), &data)
		parsedOp, err := p.parseSelect(testCase.model, data)
		if err != nil {
			t.Error(err)
		}
		diff := cmp.Diff(testCase.output, parsedOp, CmpOpts()...)
		if diff != "" {
			t.Errorf("(-want +got):\n%s", diff)
		}
	}
}
