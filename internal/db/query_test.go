package db

import (
	"fmt"
	"testing"

	"github.com/google/uuid"
)

func addTestData(db DB) {
	userId1 := uuid.MustParse("e6053eb9-3c28-4152-a4ca-c582f20fc8f0")
	userId2 := uuid.MustParse("f4cc9efe-55cf-4f05-8061-4b3b4dbc8295")
	userId3 := uuid.MustParse("d40fab5d-883b-4568-b568-b68e1cbc8292")
	postId1 := uuid.MustParse("64e3c9e2-4b4d-4009-8cb9-f8938e135926")
	postId2 := uuid.MustParse("7e374648-8a0a-4317-8768-be10f10ab743")
	postId3 := uuid.MustParse("3c5e230f-a31f-4fa5-9872-e0f6fe9c2756")
	postId4 := uuid.MustParse("eb09002a-4a7c-4fcc-86ce-b198354e6d3f")

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

	p3, err := tx.MakeRecord(Post.ID())
	if err != nil {
		panic(err)
	}
	p3.Set("id", postId3)
	p3.Set("text", "new post bout to drop")

	p4, err := tx.MakeRecord(Post.ID())
	if err != nil {
		panic(err)
	}
	p4.Set("id", postId4)
	p4.Set("text", "new post who this")

	// inserting out of order
	tx.Insert(u1)
	tx.Insert(u3)
	tx.Insert(u2)

	tx.Insert(p1)
	tx.Insert(p2)
	tx.Insert(p3)
	tx.Insert(p4)
	tx.Connect(u1.ID(), p1.ID(), UserPosts.ID())
	tx.Connect(u1.ID(), p2.ID(), UserPosts.ID())
	tx.Connect(u1.ID(), p3.ID(), UserPosts.ID())
	tx.Connect(u1.ID(), p4.ID(), UserPosts.ID())
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
	results := tx.Query(user,
		Join(post, user.Rel(userPosts)),
		Filter(post, Eq("text", "hello")),
		Aggregate(post, Some)).All()

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

	results := tx.Query(user,
		Filter(user, Eq("age", int64(32))),
		Or(user,
			tx.Subquery(Filter(user, Eq("firstName", "Andrew")),
				Join(post, user.Rel(userPosts)),
				Filter(post, Eq("text", "hello")),
				Aggregate(post, Some)),
			tx.Subquery(Filter(user, Eq("firstName", "Chase")),
				Join(post, user.Rel(userPosts)),
				Filter(post, Eq("text", "hello")),
				Aggregate(post, None)),
		),
	).All()

	if len(results) != 2 {
		t.Error("wrong number of results")
	}
}

func TestQueryCase(t *testing.T) {
	appDB := NewTest()
	AddSampleModels(appDB)
	tx := appDB.NewTx()
	relationship := tx.Ref(RelationshipInterface.ID())
	cRel := tx.Ref(ConcreteRelationshipModel.ID())
	rRel := tx.Ref(ReverseRelationshipModel.ID())
	crelRel := tx.Ref(RelationshipInterface.ID())
	rrelRef := tx.Ref(ConcreteRelationshipModel.ID())
	target, _ := tx.Schema().GetRelationshipByID(ConcreteRelationshipTarget.ID())
	referencing, _ := tx.Schema().GetRelationshipByID(ReverseRelationshipReferencing.ID())

	results := tx.Query(relationship,
		Case(relationship, cRel),
		Join(crelRel, cRel.Rel(target)),
		Case(relationship, rRel),
		Join(rrelRef, rRel.Rel(referencing)),
	).All()

	for _, res := range results {
		if res.Record.InterfaceID() == ConcreteRelationshipModel.ID() {
			_, ok := res.ToOne["target"]
			if !ok {
				err := fmt.Errorf("Didn't get target for %v\n", res)
				t.Error(err)
			}
			_, ok = res.ToOne["referencing"]
			if ok {
				err := fmt.Errorf("Got referencing for %v\n", res)
				t.Error(err)
			}
		} else if res.Record.InterfaceID() == ReverseRelationshipModel.ID() {
			_, ok := res.ToOne["target"]
			if ok {
				err := fmt.Errorf("Get target for %v\n", res)
				t.Error(err)
			}
			_, ok = res.ToOne["referencing"]
			if !ok {
				err := fmt.Errorf("Didn't get referencing for %v\n", res)
				t.Error(err)
			}
		}
	}

}

func TestLimit(t *testing.T) {
	appDB := NewTest()
	AddSampleModels(appDB)
	addTestData(appDB)
	tx := appDB.NewTx()
	users := tx.Ref(User.ID())
	results := tx.Query(users, Limit(users, 2)).All()
	if len(results) != 2 {
		err := fmt.Errorf("wrong number of results %v", results)
		t.Error(err)
	}
}

func TestLimitJoin(t *testing.T) {
	appDB := NewTest()
	AddSampleModels(appDB)
	addTestData(appDB)
	tx := appDB.NewTx()
	users := tx.Ref(User.ID())
	posts := tx.Ref(Post.ID())
	userPosts, _ := tx.Schema().GetRelationshipByID(UserPosts.ID())

	result, err := tx.Query(users,
		Filter(users, Eq("firstName", "Andrew")),
		Join(posts, users.Rel(userPosts)),
		Aggregate(posts, Include),
		Limit(posts, 2),
	).One()

	if err != nil {
		t.Fatal(err)
	}

	returnedPosts := result.GetChildRelMany(userPosts)
	if len(returnedPosts) != 2 {
		err := fmt.Errorf("wrong number of results %v", returnedPosts)
		t.Error(err)
	}

}

func TestOrder(t *testing.T) {
	expected := []string{
		"Tom", "Chase", "Andrew",
	}
	appDB := NewTest()
	AddSampleModels(appDB)
	addTestData(appDB)
	tx := appDB.NewTx()
	users := tx.Ref(User.ID())
	results := tx.Query(users, Order(users, []Sort{
		Sort{AttributeName: "firstName", Ascending: false},
	})).Records()

	for i := range results {
		if results[i].MustGet("firstName") != expected[i] {
			err := fmt.Errorf("results out of order: %v", results)
			t.Error(err)

		}
	}

}

func TestOrderJoin(t *testing.T) {
	expected := []string{
		"goodbye",
		"hello",
		"new post bout to drop",
		"new post who this",
	}
	appDB := NewTest()
	AddSampleModels(appDB)
	addTestData(appDB)
	tx := appDB.NewTx()
	users := tx.Ref(User.ID())
	posts := tx.Ref(Post.ID())
	userPosts, _ := tx.Schema().GetRelationshipByID(UserPosts.ID())

	result, err := tx.Query(users,
		Filter(users, Eq("firstName", "Andrew")),
		Join(posts, users.Rel(userPosts)),
		Aggregate(posts, Include),
		Order(posts, []Sort{Sort{AttributeName: "text", Ascending: true}}),
	).One()
	if err != nil {
		t.Fatal(err)
	}

	returnedPosts := result.GetChildRelMany(userPosts)
	for i := range returnedPosts {
		if returnedPosts[i].Record.MustGet("text") != expected[i] {
			err := fmt.Errorf("results out of order: %v", returnedPosts)
			t.Fatal(err)

		}
	}

}
