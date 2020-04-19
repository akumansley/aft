package operations

import (
	"awans.org/aft/internal/server/db"
	"github.com/json-iterator/go"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestParseFindMany(t *testing.T) {
	appDB := db.New()
	appDB.AddSampleModels()
	p := Parser{db: appDB}

	var findManyTests = []struct {
		modelName  string
		jsonString string
		output     db.FindManyOperation
	}{
		// Basic String FieldCriterion
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
		// Multiple String FieldCriterion
		{
			modelName: "user",
			jsonString: `{ 
				"firstName": "Andrew",
				"lastName": "Wansley",
				"age": 32,
			}`,
			output: db.FindManyOperation{
				ModelName: "user",
				Query: db.Query{
					FieldCriteria: []db.FieldCriterion{
						db.FieldCriterion{
							Key: "Firstname",
							Val: "Andrew",
						},
						db.FieldCriterion{
							Key: "Lastname",
							Val: "Wansley",
						},
						db.FieldCriterion{
							Key: "Age",
							Val: int64(32),
						},
					},
				},
			},
		},
		//
	}
	for _, testCase := range findManyTests {
		var data map[string]interface{}
		jsoniter.Unmarshal([]byte(testCase.jsonString), &data)
		parsedOp := p.ParseFindMany(testCase.modelName, data)
		assert.ElementsMatch(t, testCase.output.Query.FieldCriteria, parsedOp.Query.FieldCriteria)
		assert.Equal(t, testCase.output.ModelName, parsedOp.ModelName)
	}
}
