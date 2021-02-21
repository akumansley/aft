package handlers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"awans.org/aft/internal/api"
	"awans.org/aft/internal/db"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func TestUpdateManyServerParseSimple(t *testing.T) {
	appDB := db.NewTest()
	AddFunctionLiterals(appDB)
	db.AddSampleModels(appDB)

	tx := appDB.NewRWTx()
	jsonString := `{ "id": "f90e1855-dbaa-4385-9929-20efe86cccb2", "firstName":"Andrew", "lastName":"Wansley", "age": 32, "emailAddress":"andrew.wansley@gmail.com"}`
	u := api.MakeRecord(appDB.NewTx(), "user", jsonString)
	tx.Insert(u)
	jsonString2 := `{ "id": "9dd0a0c6-7e41-4107-9529-e75a5c7135cf", "firstName":"Chase", "lastName":"Hensel", "age": 32, "emailAddress":"chase.hensel@gmail.com"}`
	u2 := api.MakeRecord(appDB.NewTx(), "user", jsonString2)
	tx.Insert(u2)
	tx.Commit()

	req, err := http.NewRequest("POST", "/user.updateMany", strings.NewReader(
		`{"data":{
			"firstName":"bob"
		},
		"where": {
			"age": 32
		}
	}`))
	if err != nil {
		t.Fatal(err)
	}
	req = mux.SetURLVars(req, map[string]string{"modelName": "user", "methodName": "updateMany"})

	cs := APIHandler{DB: appDB}
	w := httptest.NewRecorder()
	err = cs.ServeHTTP(w, req)
	if err != nil {
		t.Fatal(err)
	}

	var data map[string]interface{}
	result := w.Result()
	bytes, err := ioutil.ReadAll(result.Body)
	if err != nil {
		t.Error(err)
	}
	json.Unmarshal(bytes, &data)
	assert.Equal(t, 2.0, data["count"])
}
