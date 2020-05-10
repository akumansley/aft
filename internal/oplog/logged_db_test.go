package oplog

import (
	"awans.org/aft/internal/server/db"
	"testing"
)

func TestLoggedDB(t *testing.T) {
	appDB := db.New()
	dbLog := NewMemLog()
	ldb := LoggedDB(dbLog, appDB)
	db.AddSampleModels(ldb)

	appDB2 := db.New()
	DBFromLog(appDB2, dbLog)
	appDB.DeepEquals(appDB2)
}
