package parsers

import (
	"awans.org/aft/internal/api"
	"awans.org/aft/internal/api/operations"
	"awans.org/aft/internal/db"
	"github.com/google/go-cmp/cmp"
	"github.com/json-iterator/go"
	"testing"
)

func TestParseCase(t *testing.T) {
	appDB := db.NewTest()
	db.AddSampleModels(appDB)
	tx := appDB.NewTx()
	p := Parser{Tx: tx}

	var inclusionTests = []struct {
		iface      db.Interface
		jsonString string
		output     operations.Case
	}{
		{
			db.RelationshipInterface,
			`{
			   "concreteRelationship": {
			   	"select": {"name": true},
			   },
			   "reverseRelationship": {
			   	"include": {"referencing": true},
			   }
			}`,
			operations.Case{
				Entries: []operations.CaseEntry{
					operations.CaseEntry{
						ModelID: db.ConcreteRelationshipModel.ID(),
						Select: operations.Select{
							Selecting: true,
							Fields:    api.Set{"name": {}},
						},
					},
				},
			},
		},
	}
	for _, testCase := range inclusionTests {
		var data map[string]interface{}
		jsoniter.Unmarshal([]byte(testCase.jsonString), &data)
		parsedOp, err := p.parseCase(testCase.iface, data)
		if err != nil {
			t.Fatal(err)
		}
		diff := cmp.Diff(testCase.output, parsedOp, CmpOpts()...)
		if diff != "" {
			t.Errorf("(-want +got):\n%s", diff)
		}
	}
}
