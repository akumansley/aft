package operations

import (
	"awans.org/aft/internal/server/db"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/ompluscator/dynamic-struct"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestCreateServerParseSimple(t *testing.T) {
	appDB := db.New()
	appDB.AddSampleModels()
	assert := assert.New(t)
	req, err := http.NewRequest("POST", "/user.create", strings.NewReader(
		`{"data":{
			"id":"15852d31-3bd4-4fc4-abd0-e4c7497644ab",
			"type":"user",
			"firstName":"Andrew",
			"lastName":"Wansley",
			"age": 32,
		}}`))
	if err != nil {
		t.Fatal(err)
	}
	req = mux.SetURLVars(req, map[string]string{"object": "user"})

	qs := CreateServer{DB: appDB}
	parsedReq, _ := qs.Parse(req).(CreateRequest)
	assert.IsType(parsedReq, CreateRequest{})

	op := parsedReq.Operation
	st := op.Struct

	reader := dynamicstruct.NewReader(st)
	uuid := reader.GetField("Id").Interface().(uuid.UUID)
	assert.Equal(uuid.String(), "15852d31-3bd4-4fc4-abd0-e4c7497644ab")

	firstName := reader.GetField("Firstname").String()
	assert.Equal(firstName, "Andrew")

	age := reader.GetField("Age").Int()
	assert.Equal(age, 32)
}

func TestCreateServerServe(t *testing.T) {
	appDB := db.New()
	appDB.AddSampleModels()
	jsonString := `{"id":"15852d31-3bd4-4fc4-abd0-e4c7497644ab",
					"type": "user",
					"firstName":"Andrew",
					"lastName":"Wansley", 
					"age": 32}`
	u := makeStruct(appDB, "user", jsonString)
	cOp := db.CreateOperation{
		Struct: u,
		Nested: []db.NestedOperation{},
	}
	req := CreateRequest{Operation: cOp}
	cs := CreateServer{DB: appDB}
	rr := httptest.NewRecorder()
	cs.Serve(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
	expected := `{"data":` + jsonString + `}`
	assert.JSONEq(t, expected, rr.Body.String())
}
