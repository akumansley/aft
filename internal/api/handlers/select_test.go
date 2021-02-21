package handlers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"awans.org/aft/internal/db"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func TestSelectEmpty(t *testing.T) {
	appDB := db.NewTest()
	AddFunctionLiterals(appDB)
	db.AddSampleModels(appDB)
	req, err := http.NewRequest("POST", "/user.create", strings.NewReader(
		`{"data":{
			"firstName":"Andrew",
			"lastName":"Wansley",
			"age": 32,
			"emailAddress": "andrew.wansley@gmail.com",
			"profile": {
				"create": {
					"text": "hello"
				}
			}
		},
		"select": {}

	}`))
	if err != nil {
		t.Fatal(err)
	}
	req = mux.SetURLVars(req, map[string]string{"modelName": "user", "methodName": "create"})

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

	objData := data["data"].(map[string]interface{})
	assert.Empty(t, objData)
}

func TestSelectSome(t *testing.T) {
	appDB := db.NewTest()
	AddFunctionLiterals(appDB)
	db.AddSampleModels(appDB)
	req, err := http.NewRequest("POST", "/user.create", strings.NewReader(
		`{"data":{
			"firstName":"Andrew",
			"lastName":"Wansley",
			"age": 32,
			"emailAddress": "andrew.wansley@gmail.com",
			"profile": {
				"create": {
					"text": "hello"
				}
			}
		},
		"select": {
			"firstName": true
		}

	}`))
	if err != nil {
		t.Fatal(err)
	}
	req = mux.SetURLVars(req, map[string]string{"modelName": "user", "methodName": "create"})

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

	objData := data["data"].(map[string]interface{})
	assert.Equal(t, "Andrew", objData["firstName"])
	assert.Equal(t, 1, len(objData))
}
