package api

import (
	"awans.org/aft/internal/bus"
	"awans.org/aft/internal/db"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestCreateAndrewType(t *testing.T) {
	appDB := db.New()
	eventbus := bus.New()
	db.AddSampleModels(appDB)
	req, err := http.NewRequest("POST", "/user.create", strings.NewReader(
		`{"data":{
			"firstName":"Chase",
			"lastName":"Wansley",
			"age": 32,
			"emailAddress": "andrew.wansley@gmail.com",
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

	req = mux.SetURLVars(req, map[string]string{"modelName": "user"})

	cs := CreateHandler{db: appDB, bus: eventbus}
	w := httptest.NewRecorder()
	err = cs.ServeHTTP(w, req)
	assert.Error(t, err)
}
