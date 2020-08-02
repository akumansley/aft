package parsers

import (
	"awans.org/aft/internal/db"
	"errors"
	"github.com/json-iterator/go"
	"testing"
)

func TestParseUpdateErrors(t *testing.T) {
	appDB := db.NewTest()
	db.AddSampleModels(appDB)
	tx := appDB.NewTx()
	p := Parser{Tx: tx}

	var updateErrorTests = []struct {
		modelName  string
		jsonString string
		output     error
	}{
		// Nested UpdateMany
		{
			modelName: "user",
			jsonString: `{"data": { 
			"firstName":"Chase", "profile" : true}}`,
			output: ErrInvalidStructure,
		},
	}
	for _, testCase := range updateErrorTests {
		var data map[string]interface{}
		jsoniter.Unmarshal([]byte(testCase.jsonString), &data)
		_, err := p.ParseUpdate(testCase.modelName, data)
		if !errors.Is(err, testCase.output) {
			t.Errorf("Wrong kind of error: %v, %v", testCase.output, err)
		}
	}
}
