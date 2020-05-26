package api

import (
	"awans.org/aft/internal/db"
	"github.com/stretchr/testify/assert"
	"testing"
)

func toAgeList(sts []db.Record) []int64 {
	var ages []int64
	for _, st := range sts {
		ages = append(ages, st.Get("age").(int64))
	}
	return ages
}

var testData = []string{
	`{"id":"00000000-0000-0000-0000-000000000000",
"type": "user",
"firstName":"Andrew",
"lastName":"Wansley", 
"age": 1}`,
	`{"id":"00000000-0000-0000-0000-000000000000",
"type": "user",
"firstName":"Andrew",
"lastName":"Wansley", 
"age": 2}`,
	`{"id":"00000000-0000-0000-0000-000000000000",
"type": "user",
"firstName":"Andrew",
"lastName":"Wansley", 
"age": 3}`,
}

func TestFindManyApply(t *testing.T) {
	appDB := db.New()
	db.AddSampleModels(appDB)
	tx := appDB.NewRWTx()

	// add test data
	for _, jsonString := range testData {
		st := makeRecord(tx, "user", jsonString)
		CreateOperation{Record: st}.Apply(tx)
	}
	var findManyTests = []struct {
		operation FindManyOperation
		output    []int64
	}{

		// Simple FindMany
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
			output: []int64{1, 2, 3},
		},

		// Simple FindMany
		{
			operation: FindManyOperation{
				ModelName: "user",
				Query: Query{
					FieldCriteria: []FieldCriterion{
						FieldCriterion{
							Key: "Age",
							Val: int64(3),
						},
					},
				},
			},
			output: []int64{3},
		},
	}
	for _, testCase := range findManyTests {
		result := testCase.operation.Apply(tx)
		actualAges := toAgeList(result)
		assert.ElementsMatch(t, testCase.output, actualAges)
	}
}
