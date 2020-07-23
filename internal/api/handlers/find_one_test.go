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

func TestFindOneServerParse(t *testing.T) {
	appDB := db.NewTest()
	eb := bus.New()
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

	req, err := http.NewRequest("POST", "/user.findOne", strings.NewReader(
		`{"where": {
		"firstName": "Andrew"
		}}`))
	if err != nil {
		t.Fatal(err)
	}
	req = mux.SetURLVars(req, map[string]string{"modelName": "user"})
	w := httptest.NewRecorder()
	foh := FindOneHandler{db: appDB, bus: eb}
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