package parsers

import (
	"awans.org/aft/internal/api/operations"
	"awans.org/aft/internal/db"
	"github.com/google/go-cmp/cmp"
	"github.com/json-iterator/go"
	"testing"
)

func TestParseInclude(t *testing.T) {
	appDB := db.NewTest()
	db.AddSampleModels(appDB)
	p := Parser{Tx: appDB.NewTx()}

	var inclusionTests = []struct {
		modelName  string
		jsonString string
		output     operations.Include
	}{
		// Simple Include
		{
			modelName: "user",
			jsonString: `{ 
			   "profile": true
			}`,
			output: operations.Include{
				Includes: []operations.Inclusion{
					operations.Inclusion{
						Relationship: db.UserProfile,
						Where:        operations.Where{},
					},
				},
			},
		},
	}
	for _, testCase := range inclusionTests {
		var data map[string]interface{}
		jsoniter.Unmarshal([]byte(testCase.jsonString), &data)
		parsedOp, err := p.ParseInclude(testCase.modelName, data)
		if err != nil {
			t.Error(err)
		}
		diff := cmp.Diff(testCase.output, parsedOp, operations.CmpOpts()...)
		if diff != "" {
			t.Errorf("(-want +got):\n%s", diff)
		}
	}
}
