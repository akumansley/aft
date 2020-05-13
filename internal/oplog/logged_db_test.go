package oplog

import (
	"awans.org/aft/internal/db"
	"awans.org/aft/internal/model"
	"encoding/json"
	"testing"
)

func makeRecord(tx db.Tx, modelName string, jsonValue string) model.Record {
	st := tx.MakeRecord(modelName)
	json.Unmarshal([]byte(jsonValue), &st)
	return st
}

func TestLoggedDB(t *testing.T) {
	appDB := db.New()
	db.AddSampleModels(appDB)
	jsonString := `{ "firstName":"Andrew", "lastName":"Wansley", "age": 32}`
	u := makeRecord(appDB.NewTx(), "user", jsonString)

	dbLog := NewMemLog()
	ldb := LoggedDB(dbLog, appDB)
	rwtx := ldb.NewRWTx()
	rwtx.Insert(u)
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

	dbLog, err := OpenGobLog("/Users/awans/Desktop/test.dbl")
	if err != nil {
		t.Fatal(err)
	}
	defer dbLog.Close()

	ldb := LoggedDB(dbLog, appDB)
	rwtx := ldb.NewRWTx()
	rwtx.Insert(u)
	rwtx.Commit()

	appDB2 := db.New()
	db.AddSampleModels(appDB2)
	DBFromLog(appDB2, dbLog)
	eq := appDB.DeepEquals(appDB2)
	if !eq {
		t.Errorf("Not equal")
	}
}
