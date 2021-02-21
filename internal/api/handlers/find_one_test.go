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

func TestFindOneServerParse(t *testing.T) {
	appDB := db.NewTest()
	AddFunctionLiterals(appDB)
	db.AddSampleModels(appDB)

	tx := appDB.NewRWTx()
	jsonString := `{ "firstName":"Andrew", "lastName":"Wansley", "age": 32, "emailAddress":"andrew.wansley@gmail.com"}`
	u := api.MakeRecord(appDB.NewTx(), "user", jsonString)
	tx.Insert(u)
	tx.Commit()

	req, err := http.NewRequest("POST", "/user.findOne", strings.NewReader(
		`{"where": {
		"firstName": "Andrew"
		}}`))
	if err != nil {
		t.Fatal(err)
	}
	req = mux.SetURLVars(req, map[string]string{"modelName": "user", "methodName": "findOne"})
	w := httptest.NewRecorder()
	foh := APIHandler{DB: appDB}
	err = foh.ServeHTTP(w, req)
	if err != nil {
		t.Fatal(err)
	}

	result := w.Result()
	bytes, err := ioutil.ReadAll(result.Body)
	if err != nil {
		t.Fatal(err)
	}

	var data map[string]interface{}
	json.Unmarshal(bytes, &data)
	objData := data["data"].(map[string]interface{})
	assert.Equal(t, "Andrew", objData["firstName"])
	assert.Equal(t, "Wansley", objData["lastName"])
	assert.Equal(t, "andrew.wansley@gmail.com", objData["emailAddress"])
	assert.Equal(t, 32.0, objData["age"])
}
