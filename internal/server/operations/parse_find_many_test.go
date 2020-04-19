package operations

import (
	"awans.org/aft/internal/server/db"
	"github.com/go-test/deep"
	"github.com/google/uuid"
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
			}`,
			output: db.FindOneOperation{
				ModelName: "user",
				UniqueQuery: db.UniqueQuery{
					Key: "Id",
					Val: uuid.MustParse("15852d31-3bd4-4fc4-abd0-e4c7497644ab"),
				},
			},
		},
	}
	for _, testCase := range findOneTests {
		var data map[string]interface{}
		jsoniter.Unmarshal([]byte(testCase.jsonString), &data)
		parsedOp := p.ParseFindOne(testCase.modelName, data)
		if diff := deep.Equal(parsedOp, testCase.output); diff != nil {
			t.Error(diff)
		}
	}
}
