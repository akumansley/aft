package parsers

import (
	"awans.org/aft/internal/db"
	"errors"
	"github.com/json-iterator/go"
	"testing"
)

func TestParseCreateErrors(t *testing.T) {
	appDB := db.NewTest()
	db.AddSampleModels(appDB)
	p := Parser{Tx: appDB.NewTx()}

	var createErrorTests = []struct {
		modelName  string
		jsonString string
		output     error
	}{
		// Simple Create
		{
			modelName: "user",
			jsonString: `{ 
			"id":"15852d31-3bd4-4fc4-abd0-e4c7497644ab",
			"firstNamez":"Andrew",
			"lastName":"Wansley",
			"emailAddress":"andrew.wansley@gmail.com",
			"age": 32}`,
			output: ErrUnusedKeys,
		},
	}
	for _, testCase := range createErrorTests {
		var data map[string]interface{}
		jsoniter.Unmarshal([]byte(testCase.jsonString), &data)
		_, err := p.ParseCreate(testCase.modelName, data)
		if !errors.Is(err, testCase.output) {
			t.Errorf("Wrong kind of error: %v, %v", testCase.output, err)
		}
	}
}
