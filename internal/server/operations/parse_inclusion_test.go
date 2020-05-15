package operations

import (
	"awans.org/aft/internal/db"
	"github.com/go-test/deep"
	"github.com/json-iterator/go"
	"testing"
)

func TestParseInclude(t *testing.T) {
	appDB := db.New()
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
						Relationship: db.User.Relationships["profile"],
						Query:        Query{},
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
		if diff := deep.Equal(parsedOp, testCase.output); diff != nil {
			t.Error(diff)
		}
	}
}
