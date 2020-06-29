package db

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"testing"
)

func addTestData(db DB) {
	userId1 := uuid.MustParse("e6053eb9-3c28-4152-a4ca-c582f20fc8f0")
	userId2 := uuid.MustParse("f4cc9efe-55cf-4f05-8061-4b3b4dbc8295")
	userId3 := uuid.MustParse("d40fab5d-883b-4568-b568-b68e1cbc8292")
	postId1 := uuid.MustParse("64e3c9e2-4b4d-4009-8cb9-f8938e135926")
	postId2 := uuid.MustParse("7e374648-8a0a-4317-8768-be10f10ab743")

	tx := db.NewRWTx()
	u1, err := tx.MakeRecord(User.ID())
	if err != nil {
		panic(err)
	}

	u1.Set("id", userId1)
	u1.Set("firstName", "Andrew")
	u1.Set("age", int64(32))

	u2, err := tx.MakeRecord(User.ID())
	if err != nil {
		panic(err)
	}
	u2.Set("id", userId2)
	u2.Set("firstName", "Chase")
	u2.Set("age", int64(32))

	u3, err := tx.MakeRecord(User.ID())
	if err != nil {
		panic(err)
	}
	u3.Set("id", userId3)
	u3.Set("firstName", "Tom")
	u3.Set("age", int64(32))

	p1, err := tx.MakeRecord(Post.ID())
	if err != nil {
		panic(err)
	}
	p1.Set("id", postId1)
	p1.Set("text", "hello")

	p2, err := tx.MakeRecord(Post.ID())
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
	tx.Connect(u1.ID(), p1.ID(), UserPosts.ID())
	tx.Connect(u1.ID(), p2.ID(), UserPosts.ID())
	tx.Commit()
}

func TestQueryJoinMany(t *testing.T) {
	appDB := NewTest()
	AddSampleModels(appDB)
	addTestData(appDB)
	tx := appDB.NewTx()
	user := tx.Ref(User.ID())
	post := tx.Ref(Post.ID())
	userPosts, err := tx.Schema().GetRelationshipByID(UserPosts.ID())
	if err != nil {
		t.Fatal(err)
	}
	results := tx.Query(user).Join(post, user.Rel(userPosts)).Filter(post, Eq("text", "hello")).Aggregate(post, Some).All()
	if len(results) != 1 {
		t.Error("wrong number of results")
	}
}

func TestQueryOr(t *testing.T) {
	appDB := NewTest()
	AddSampleModels(appDB)
	addTestData(appDB)
	tx := appDB.NewTx()
	user := tx.Ref(User.ID())
	post := tx.Ref(Post.ID())
	userPosts, _ := tx.Schema().GetRelationshipByID(UserPosts.ID())

	results := tx.Query(user).Filter(user, Eq("age", int64(32))).Or(user,
		Filter(user, Eq("firstName", "Andrew")).Join(post, user.Rel(userPosts)).Filter(post, Eq("text", "hello")).Aggregate(post, Some),
		Filter(user, Eq("firstName", "Chase")).Join(post, user.Rel(userPosts)).Filter(post, Eq("text", "hello")).Aggregate(post, None),
	).All()

	bytes, _ := json.Marshal(results)
	fmt.Printf("results: %v\n", string(bytes))
	if len(results) != 2 {
		t.Error("wrong number of results")
	}
}
