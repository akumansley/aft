package operations

import (
	"awans.org/aft/internal/server/db"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	db.InitDB()
	AddSampleModels()
	os.Exit(m.Run())
}

func AddSampleModels() {
	db.DB.Insert(db.User.Name, &db.User)
	db.DB.Insert(db.Profile.Name, &db.Profile)
	db.DB.Insert(db.Post.Name, &db.Post)
}
