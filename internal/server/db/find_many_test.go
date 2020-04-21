package db

import (
	"github.com/google/uuid"
	"github.com/ompluscator/dynamic-struct"
	"github.com/stretchr/testify/assert"
	"testing"
)

func getIdAsString(st interface{}) string {
	reader := dynamicstruct.NewReader(st)
	id := reader.GetField("Id").Interface().(uuid.UUID)
	return id.String()
}

func toIdList(sts []interface{}) []string {
	var ids []string
	for _, st := range sts {
		ids = append(ids, getIdAsString(st))
	}
	return ids
}

var testData = []string{
	`{"id":"15852d31-3bd4-4fc4-abd0-e4c7497644ab",
"type": "user",
"firstName":"Andrew",
"lastName":"Wansley", 
"age": 32}`,
	`{"id":"9514ca6b-ef2e-4b7f-8cb2-5aa1557f5ea1",
"type": "user",
"firstName":"Andrew",
"lastName":"Wansley", 
"age": 12}`,
	`{"id":"51328560-edec-4a0d-a475-0d9a76b09103",
"type": "user",
"firstName":"Andrew",
"lastName":"Wansley", 
"age": 16}`,
}

func TestFindManyApply(t *testing.T) {
	appDB := New()
	appDB.AddSampleModels()

	// add test data
	for _, jsonString := range testData {
		st := makeStruct(appDB, "user", jsonString)
		CreateOperation{Struct: st}.Apply(appDB)
	}
	var findManyTests = []struct {
		operation FindManyOperation
		output    []string
	}{
		{
			operation: FindManyOperation{
				ModelName: "user",
				Query: Query{
					FieldCriteria: []FieldCriterion{
						FieldCriterion{
							Key: "Firstname",
							Val: "Andrew",
						},
					},
				},
			},
			output: []string{
				"15852d31-3bd4-4fc4-abd0-e4c7497644ab",
				"9514ca6b-ef2e-4b7f-8cb2-5aa1557f5ea1",
				"51328560-edec-4a0d-a475-0d9a76b09103",
			},
		},
	}
	for _, testCase := range findManyTests {
		result := testCase.operation.Apply(appDB)
		rList := result.([]interface{})
		actualIds := toIdList(rList)
		assert.ElementsMatch(t, testCase.output, actualIds)
	}
}
