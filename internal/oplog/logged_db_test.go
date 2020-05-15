package oplog

import (
	"awans.org/aft/internal/db"
	"encoding/json"
	"io/ioutil"
	"os"
	"testing"
)

func makeRecord(tx db.Tx, modelName string, jsonValue string) db.Record {
	st := tx.MakeRecord(modelName)
	json.Unmarshal([]byte(jsonValue), &st)
	return st
}

func TestLoggedDB(t *testing.T) {
	appDB := db.New()
	db.AddSampleModels(appDB)
	jsonString := `{ "firstName":"Andrew", "lastName":"Wansley", "age": 32}`
	u := makeRecord(appDB.NewTx(), "user", jsonString)
	jsonString = `{ "text":"hello.." }`
	p := makeRecord(appDB.NewTx(), "profile", jsonString)

	dbLog := NewMemLog()
	ldb := LoggedDB(dbLog, appDB)
	rwtx := ldb.NewRWTx()
	rwtx.Insert(u)
	rwtx.Insert(p)
	rwtx.Connect(u, p, db.User.Relationships["profile"])
	rwtx.Commit()

	appDB2 := db.New()
	db.AddSampleModels(appDB2)
	DBFromLog(appDB2, dbLog)
	eq := appDB.DeepEquals(appDB2)
	if !eq {
		t.Errorf("Not equal")
	}
}

func TestGobLoggedDB(t *testing.T) {
	appDB := db.New()
	db.AddSampleModels(appDB)
	jsonString := `{ "firstName":"Andrew", "lastName":"Wansley", "age": 32}`
	u := makeRecord(appDB.NewTx(), "user", jsonString)
	jsonString = `{ "text":"hello.." }`
	p := makeRecord(appDB.NewTx(), "profile", jsonString)

	tmpFile, err := ioutil.TempFile(os.TempDir(), "aft-test-")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpFile.Name())

	dbLog, err := OpenGobLog(tmpFile.Name())
	if err != nil {
		t.Fatal(err)
	}
	defer dbLog.Close()

	ldb := LoggedDB(dbLog, appDB)
	rwtx := ldb.NewRWTx()
	rwtx.Insert(u)
	rwtx.Insert(p)
	rwtx.Connect(u, p, db.User.Relationships["profile"])
	rwtx.Commit()

	appDB2 := db.New()
	db.AddSampleModels(appDB2)
	DBFromLog(appDB2, dbLog)
	eq := appDB.DeepEquals(appDB2)
	if !eq {
		t.Errorf("Not equal")
	}
}
