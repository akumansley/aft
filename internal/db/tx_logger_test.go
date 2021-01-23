package db

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"testing"

	"awans.org/aft/internal/bus"
	"awans.org/aft/internal/oplog"
	"github.com/google/uuid"
)

func makeRecord(tx Tx, modelName string, jsonValue string) Record {
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
	appDB := New(b)
	AddSampleModels(appDB)
	jsonString := `{ 
		"id":"b5833b1e-4f78-486f-b05f-c8997cdf5cae", 
		"firstName":"Andrew", 
		"lastName":"Wansley", 
		"age": 32}`
	u := makeRecord(appDB.NewTx(), "user", jsonString)
	jsonString = `{ 
		"id":"45380fae-982d-4573-8734-7acacb279bfe", 
		"text":"hello.." }`
	p := makeRecord(appDB.NewTx(), "profile", jsonString)

	n := u.DeepCopy()
	n.Set("firstName", "Chase")
	n.Set("id", uuid.MustParse("cb56658a-9525-4102-a9be-e2645d6c9975"))
	dbLog := DBOpLog(appDB.Builder(), oplog.NewMemLog())
	txLogger := MakeTransactionLogger(dbLog)
	b.RegisterHandler(txLogger)
	rwtx := appDB.NewRWTx()
	rwtx.Insert(u)
	rwtx.Insert(p)
	rwtx.Connect(u.ID(), p.ID(), UserProfile.ID())
	rwtx.Update(u, n)
	rwtx.Commit()

	appDB2 := NewTest()
	AddSampleModels(appDB2)
	DBFromLog(appDB2, dbLog)
	eq := appDB.DeepEquals(appDB2)
	if !eq {
		t.Errorf("Not equal")
	}
}

func TestGobLoggedDB(t *testing.T) {
	b := bus.New()
	appDB := New(b)
	AddSampleModels(appDB)
	jsonString := `{ 
		"id":"69bbe86b-03cd-4841-8673-806e6f965015",
		"firstName":"Andrew", 
		"lastName":"Wansley", 
		"age": 32}`
	u := makeRecord(appDB.NewTx(), "user", jsonString)
	jsonString = `{ 
		"id": "efda33fa-5196-4433-924e-8f14005fe5a9",
		"text": "hello.." }`
	p := makeRecord(appDB.NewTx(), "profile", jsonString)

	tmpFile, err := ioutil.TempFile(os.TempDir(), "aft-test-")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpFile.Name())

	logStore, err := oplog.OpenDiskLog(tmpFile.Name())
	if err != nil {
		t.Fatal(err)
	}
	dbLog := DBOpLog(appDB.Builder(), logStore)
	defer dbLog.Close()

	n := u.DeepCopy()
	n.Set("firstName", "Chase")
	n.Set("id", uuid.MustParse("768bb99c-4cc5-452f-9ba5-d19c6bc5dc7c"))
	txLogger := MakeTransactionLogger(dbLog)
	b.RegisterHandler(txLogger)
	rwtx := appDB.NewRWTx()
	rwtx.Insert(u)
	rwtx.Insert(p)
	rwtx.Connect(u.ID(), p.ID(), UserProfile.ID())
	rwtx.Update(u, n)
	rwtx.Commit()

	appDB2 := NewTest()
	AddSampleModels(appDB2)
	DBFromLog(appDB2, dbLog)
	eq := appDB.DeepEquals(appDB2)
	if !eq {
		t.Errorf("Not equal")
	}
}
