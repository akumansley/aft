package api

import (
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

func TestCreateServerParseSimple(t *testing.T) {
	appDB := db.New()
	eventbus := bus.New()
	db.AddSampleModels(appDB)
	req, err := http.NewRequest("POST", "/user.create", strings.NewReader(
		`{"data":{
			"firstName":"Andrew",
			"lastName":"Wansley",
			"age": 32,
			"email": "andrew.wansley@gmail.com",
			"profile": {
				"create": {
					"text": "hello"
				}
			}
		},
		"include": {
			"profile": true
		}

	}`))
	if err != nil {
		t.Fatal(err)
	}
	req = mux.SetURLVars(req, map[string]string{"modelName": "user"})

	cs := CreateHandler{db: appDB, bus: eventbus}
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
	assert.Equal(t, "Wansley", objData["lastName"])
	assert.Equal(t, 32.0, objData["age"])
	profileData := objData["profile"].(map[string]interface{})
	assert.Equal(t, "hello", profileData["text"])

}
