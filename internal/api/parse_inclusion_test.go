package api

import (
	"awans.org/aft/internal/db"
	"github.com/google/go-cmp/cmp"
	"github.com/json-iterator/go"
	"testing"
)

func TestParseInclude(t *testing.T) {
	appDB := db.NewTest()
	db.AddSampleModels(appDB)
	p := Parser{tx: appDB.NewTx()}

	var inclusionTests = []struct {
		modelName  string
		jsonString string
		output     Include
	}{
		// Simple Include
		{
			modelName: "user",
			jsonString: `{ 
			   "profile": true
			}`,
			output: Include{
				Includes: []Inclusion{
					Inclusion{
						Relationship: db.UserProfile,
						Where:        Where{},
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
		diff := cmp.Diff(testCase.output, parsedOp, CmpOpts()...)
		if diff != "" {
			t.Errorf("(-want +got):\n%s", diff)
		}
	}
}
