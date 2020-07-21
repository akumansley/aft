package operations

import (
	"awans.org/aft/internal/db"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"testing"
)

var (
	userId1 = uuid.MustParse("e6053eb9-3c28-4152-a4ca-c582f20fc8f0")
	userId2 = uuid.MustParse("f4cc9efe-55cf-4f05-8061-4b3b4dbc8295")
	userId3 = uuid.MustParse("d40fab5d-883b-4568-b568-b68e1cbc8292")
	postId1 = uuid.MustParse("64e3c9e2-4b4d-4009-8cb9-f8938e135926")
	postId2 = uuid.MustParse("7e374648-8a0a-4317-8768-be10f10ab743")
)

func addTestData(appDB db.DB) {
	tx := appDB.NewRWTx()
	u1, err := tx.MakeRecord(db.User.ID())
	if err != nil {
		panic(err)
	}
	u1.Set("id", userId1)
	u1.Set("firstName", "Gid")
	u1.Set("age", int64(4))

	u2, err := tx.MakeRecord(db.User.ID())
	if err != nil {
		panic(err)
	}
	u2.Set("id", userId2)
	u2.Set("firstName", "Chase")
	u2.Set("age", int64(5))

	u3, err := tx.MakeRecord(db.User.ID())
	if err != nil {
		panic(err)
	}
	u3.Set("id", userId3)
	u3.Set("firstName", "Tom")
	u3.Set("age", int64(6))

	p1, err := tx.MakeRecord(db.Post.ID())
	if err != nil {
		panic(err)
	}
	p1.Set("id", postId1)
	p1.Set("text", "hello")

	p2, err := tx.MakeRecord(db.Post.ID())
	if err != nil {
		panic(err)
	}
	p2.Set("id", postId2)
	p2.Set("text", "goodbye")

	tx.Insert(u1)
	tx.Insert(u2)
	tx.Insert(u3)
	tx.Insert(p1)
	tx.Insert(p2)
	tx.Connect(u1.ID(), p1.ID(), db.UserPosts.ID())
	tx.Connect(u1.ID(), p2.ID(), db.UserPosts.ID())

	tx.Commit()
}

func toAgeList(sts []db.Record) []int64 {
	var ages []int64
	for _, st := range sts {
		age, _ := st.Get("age")
		ages = append(ages, age.(int64))
	}
	return ages
}

var testData = []string{
	`{"id":"64cf3e64-2cad-4065-9f99-23266eb27632",
"type": "user",
"firstName":"Andrew",
"lastName":"Wansley", 
"age": 1}`,
	`{"id":"d193faaa-0a18-47d7-a3f0-ef3fd2e62a0a",
"type": "user",
"firstName":"Andrew",
"lastName":"Wansley", 
"age": 2}`,
	`{"id":"14f8e532-d918-465d-99c1-dade3413dfa7",
"type": "user",
"firstName":"Andrew",
"lastName":"Wansley", 
"age": 3}`,
}

func TestFindManyApply(t *testing.T) {
	appDB := db.NewTest()
	db.AddSampleModels(appDB)
	addTestData(appDB)
	tx := appDB.NewRWTx()

	// add test data
	for _, jsonString := range testData {
		st := MakeRecord(tx, "user", jsonString)
		CreateOperation{Record: st}.Apply(tx)
	}
	var findManyTests = []struct {
		operation FindManyOperation
		output    []int64
	}{

		// Simple FindMany
		{
			operation: FindManyOperation{
				ModelID: db.User.ID(),
				Where: Where{
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

		// FindMany with aggregate Related
		{
			operation: FindManyOperation{
				ModelID: db.User.ID(),
				Where: Where{
					AggregateRelationshipCriteria: []AggregateRelationshipCriterion{
						AggregateRelationshipCriterion{
							RelationshipCriterion: RelationshipCriterion{
								Relationship: db.UserPosts,
								Where: Where{
									FieldCriteria: []FieldCriterion{
										FieldCriterion{
											Key: "text",
											Val: "hello",
										},
									},
								},
							},
							Aggregation: db.Some,
						},
					},
				},
			},
			output: []int64{4},
		},

		// FindMany with OR
		{
			operation: FindManyOperation{
				ModelID: db.User.ID(),
				Where: Where{
					Or: []Where{

						Where{
							AggregateRelationshipCriteria: []AggregateRelationshipCriterion{
								AggregateRelationshipCriterion{
									RelationshipCriterion: RelationshipCriterion{
										Relationship: db.UserPosts,
										Where: Where{
											FieldCriteria: []FieldCriterion{
												FieldCriterion{
													Key: "text",
													Val: "goodbye",
												},
											},
										},
									},
									Aggregation: db.Some,
								},
							},
						},

						Where{
							AggregateRelationshipCriteria: []AggregateRelationshipCriterion{
								AggregateRelationshipCriterion{
									RelationshipCriterion: RelationshipCriterion{
										Relationship: db.UserPosts,
										Where: Where{
											FieldCriteria: []FieldCriterion{
												FieldCriterion{
													Key: "text",
													Val: "hello",
												},
											},
										},
									},
									Aggregation: db.Some,
								},
							},
						},
					},
				},
			},
			output: []int64{4},
		},
	}
	for _, testCase := range findManyTests {
		result := testCase.operation.Apply(tx)
		actualAges := toAgeList(result)
		assert.ElementsMatch(t, testCase.output, actualAges)
	}
}
