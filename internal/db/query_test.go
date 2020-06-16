package db

import (
	"encoding/json"
	"fmt"
	"testing"
)

const (
	userId1 = "e6053eb9-3c28-4152-a4ca-c582f20fc8f0"
	userId2 = "f4cc9efe-55cf-4f05-8061-4b3b4dbc8295"
	postId1 = "64e3c9e2-4b4d-4009-8cb9-f8938e135926"
	postId2 = "7e374648-8a0a-4317-8768-be10f10ab743"
)

func addTestData(db DB) {
	tx := db.NewRWTx()
	u1 := tx.MakeRecord(User.ID)
	u1.Set("id", userId1)

	u2 := tx.MakeRecord(User.ID)
	u2.Set("id", userId2)

	p1 := tx.MakeRecord(Post.ID)
	p1.Set("id", postId1)
	p1.Set("text", "hello")

	p2 := tx.MakeRecord(Post.ID)
	p2.Set("id", postId2)
	p2.Set("text", "goodbye")

	tx.Insert(u1)
	tx.Insert(u2)
	tx.Insert(p1)
	tx.Insert(p2)
	tx.Connect(u1, p1, UserPosts)
	tx.Connect(u1, p2, UserPosts)

	tx.Commit()
}

func TestQuery(t *testing.T) {
	appDB := New()
	AddSampleModels(appDB)
	addTestData(appDB)
	tx := appDB.NewTx()
	user := tx.Ref(User.ID)
	post := tx.Ref(Post.ID)
	results := tx.Query(user).Join(post, user.Rel("posts")).Filter(post, Eq("text", "hello")).All()
	bytes, _ := json.Marshal(results)
	fmt.Printf("results: %v", string(bytes))
	t.Error(results)

}
