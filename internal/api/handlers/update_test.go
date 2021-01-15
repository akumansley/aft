package handlers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"awans.org/aft/internal/api"
	"awans.org/aft/internal/bus"
	"awans.org/aft/internal/db"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func TestUpdateServerParseSimple(t *testing.T) {
	appDB := db.NewTest()
	eventbus := bus.New()
	db.AddSampleModels(appDB)

	tx := appDB.NewRWTx()
	jsonString := `{"firstName":"Andrew", "lastName":"Wansley", "age": 32, "emailAddress":"andrew.wansley@gmail.com"}`
	u := api.MakeRecord(appDB.NewTx(), "user", jsonString)
	tx.Insert(u)
	tx.Commit()

	req, err := http.NewRequest("POST", "/user.update", strings.NewReader(
		`{
		"data":{
			"firstName":"Chase"
		},
		"where": {
			"firstName": "Andrew"
		},
		"select" : {
			"firstName" : true,
			"lastName" : true,
			"profile" : true
		}
	}`))
	if err != nil {
		t.Fatal(err)
	}
	req = mux.SetURLVars(req, map[string]string{"modelName": "user"})

	cs := UpdateHandler{db: appDB, bus: eventbus}
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

	objData := data["data"].(map[string]interface{})
	assert.Equal(t, "Chase", objData["firstName"])
	assert.Equal(t, "Wansley", objData["lastName"])
	_, ok := objData["age"]
	assert.Equal(t, false, ok)
}
