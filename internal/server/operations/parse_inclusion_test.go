package operations

import (
	"awans.org/aft/internal/server/db"
	"github.com/go-test/deep"
	"github.com/json-iterator/go"
	"testing"
)

func TestParseInclude(t *testing.T) {
	appDB := db.New()
	appDB.AddSampleModels()
	p := Parser{db: appDB}

	var inclusionTests = []struct {
		modelName  string
		jsonString string
		output     db.Include
	}{
		// Simple Include
		{
			modelName: "user",
			jsonString: `{ 
			   "profile": true
			}`,
			output: db.Include{
				Includes: []db.Inclusion{
					db.Inclusion{
						Relationship: db.User.Relationships["profile"],
						Query:        db.Query{},
					},
				},
			},
		},
	}
	for _, testCase := range inclusionTests {
		var data map[string]interface{}
		jsoniter.Unmarshal([]byte(testCase.jsonString), &data)
		parsedOp := p.ParseInclude(testCase.modelName, data)
		if diff := deep.Equal(parsedOp, testCase.output); diff != nil {
			t.Error(diff)
		}
	}
}
