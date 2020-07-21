package handlers

import (
	"awans.org/aft/internal/api/operations"
	"awans.org/aft/internal/bus"
	"awans.org/aft/internal/db"
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestUpdateServerParseSimple(t *testing.T) {
	appDB := db.NewTest()
	eventbus := bus.New()
	db.AddSampleModels(appDB)

	jsonString := `{ "firstName":"Andrew", "lastName":"Wansley", "age": 32, "emailAddress":"andrew.wansley@gmail.com"}`
	u := operations.MakeRecord(appDB.NewTx(), "user", jsonString)
	cOp := operations.CreateOperation{
		Record: u,
		Nested: []operations.NestedOperation{},
	}
	tx := appDB.NewRWTx()
	cOp.Apply(tx)
	tx.Commit()

	req, err := http.NewRequest("POST", "/user.update", strings.NewReader(
		`{"data":{
			"firstName":"Chase"
		},
		"where": {
			"firstName": "Andrew"
		}
	}`))
	if err != nil {
		t.Fatal(err)
	}
	req = mux.SetURLVars(req, map[string]string{"modelName": "user"})

	cs := UpdateHandler{DB: appDB, Bus: eventbus}
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
}
