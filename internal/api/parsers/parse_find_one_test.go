package parsers

import (
	"awans.org/aft/internal/api/operations"
	"awans.org/aft/internal/db"
	"github.com/go-test/deep"
	"github.com/google/uuid"
	"github.com/json-iterator/go"
	"testing"
)

func TestParseFindOne(t *testing.T) {
	appDB := db.NewTest()
	db.AddSampleModels(appDB)
	p := Parser{Tx: appDB.NewTx()}

	var findOneTests = []struct {
		modelName  string
		jsonString string
		output     interface{}
	}{
		// Simple FindOne
		{
			modelName: "user",
			jsonString: `{ 
			"id":"15852d31-3bd4-4fc4-abd0-e4c7497644ab",
			}`,
			output: operations.FindOneOperation{
				ModelID: db.User.ID(),
				UniqueQuery: operations.UniqueQuery{
					Key: "Id",
					Val: uuid.MustParse("15852d31-3bd4-4fc4-abd0-e4c7497644ab"),
				},
			},
		},
	}
	for _, testCase := range findOneTests {
		var data map[string]interface{}
		jsoniter.Unmarshal([]byte(testCase.jsonString), &data)
		parsedOp, err := p.ParseFindOne(testCase.modelName, data)
		if err != nil {
			t.Error(err)
		}
		if diff := deep.Equal(parsedOp, testCase.output); diff != nil {
			t.Error(diff)
		}
	}
}