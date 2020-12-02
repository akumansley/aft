package oplog

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"testing"

	"awans.org/aft/internal/bus"
	"awans.org/aft/internal/db"
)

func makeRecord(tx db.Tx, modelName string, jsonValue string) db.Record {
	m, _ := tx.Schema().GetModel(modelName)
	st, err := tx.MakeRecord(m.ID())
	if err != nil {
		panic(err)
	}
	json.Unmarshal([]byte(jsonValue), &st)
	return st
}

func TestTxLogger(t *testing.T) {
	b := bus.New()
	appDB := db.New(b)
	db.AddSampleModels(appDB)
	jsonString := `{ "firstName":"Andrew", "lastName":"Wansley", "age": 32}`
	u := makeRecord(appDB.NewTx(), "user", jsonString)
	jsonString = `{ "text":"hello.." }`
	p := makeRecord(appDB.NewTx(), "profile", jsonString)

	n := u.DeepCopy()
	n.Set("firstName", "Chase")
	dbLog := NewMemLog()
	txLogger := MakeTransactionLogger(dbLog)
	b.RegisterHandler(txLogger)
	rwtx := appDB.NewRWTx()
	rwtx.Insert(u)
	rwtx.Insert(p)
	rwtx.Connect(u.ID(), p.ID(), db.UserProfile.ID())
	rwtx.Update(u, n)
	rwtx.Commit()

	appDB2 := db.NewTest()
	db.AddSampleModels(appDB2)
	DBFromLog(appDB2, dbLog)
	eq := appDB.DeepEquals(appDB2)
	if !eq {
		t.Errorf("Not equal")
	}
}

func TestGobLoggedDB(t *testing.T) {
	b := bus.New()
	appDB := db.New(b)
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

	n := u.DeepCopy()
	n.Set("firstName", "Chase")
	txLogger := MakeTransactionLogger(dbLog)
	b.RegisterHandler(txLogger)
	rwtx := appDB.NewRWTx()
	rwtx.Insert(u)
	rwtx.Insert(p)
	rwtx.Connect(u.ID(), p.ID(), db.UserProfile.ID())
	rwtx.Update(u, n)
	rwtx.Commit()

	appDB2 := db.NewTest()
	db.AddSampleModels(appDB2)
	DBFromLog(appDB2, dbLog)
	eq := appDB.DeepEquals(appDB2)
	if !eq {
		t.Errorf("Not equal")
	}
}
