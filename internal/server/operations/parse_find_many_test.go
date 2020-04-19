package operations

import (
	"awans.org/aft/internal/server/db"
	"github.com/go-test/deep"
	"github.com/json-iterator/go"
	"testing"
)

func TestParseFindMany(t *testing.T) {
	appDB := db.New()
	appDB.AddSampleModels()
	p := Parser{db: appDB}

	var findManyTests = []struct {
		modelName  string
		jsonString string
		output     interface{}
	}{
		// Simple FindMany
		{
			modelName: "user",
			jsonString: `{ 
				"firstName": "Andrew"
			}`,
			output: db.FindManyOperation{
				ModelName: "user",
				Query: db.Query{
					FieldCriteria: []db.FieldCriterion{
						db.FieldCriterion{
							Key: "Firstname",
							Val: "Andrew",
						},
					},
				},
			},
		},
	}
	for _, testCase := range findManyTests {
		var data map[string]interface{}
		jsoniter.Unmarshal([]byte(testCase.jsonString), &data)
		parsedOp := p.ParseFindMany(testCase.modelName, data)
		if diff := deep.Equal(parsedOp, testCase.output); diff != nil {
			t.Error(diff)
		}
	}
}
