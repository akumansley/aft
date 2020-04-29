package operations

import (
	"awans.org/aft/internal/server/db"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/ompluscator/dynamic-struct"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func getId(st interface{}) uuid.UUID {
	reader := dynamicstruct.NewReader(st)
	id := reader.GetField("Id").Interface().(uuid.UUID)
	return id
}

func TestFindOneServerParse(t *testing.T) {
	appDB := db.New()
	appDB.AddSampleModels()
	req, err := http.NewRequest("POST", "/user.query", strings.NewReader(
		`{"where": {
		"id": "2b1e9f08-38a9-4a36-b653-0fa0cbc8cad2"
		}}`))
	if err != nil {
		t.Fatal(err)
	}
	req = mux.SetURLVars(req, map[string]string{"object": "user"})
	fos := FindOneServer{DB: appDB}
	_, ok := fos.Parse(req).(FindOneRequest)
	if !ok {
		t.Fatal("Didn't return a FindOneRequest")
	}
}

func TestFindOneServerServe(t *testing.T) {
	appDB := db.New()
	appDB.AddSampleModels()
	jsonString := `{ "firstName":"Andrew", "lastName":"Wansley", "age": 32}`
	u := makeStruct(appDB, "user", jsonString)
	cOp := db.CreateOperation{
		Struct: u,
		Nested: []db.NestedOperation{},
	}
	cOp.Apply(appDB)

	req := FindOneRequest{Operation: db.FindOneOperation{
		ModelName: "user",
		UniqueQuery: db.UniqueQuery{
			Key: "Id",
			Val: getId(u),
		},
	}}
	fos := FindOneServer{DB: appDB}
	rr := httptest.NewRecorder()
	fos.Serve(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
	var data map[string]interface{}
	json.Unmarshal(rr.Body.Bytes(), &data)
	objData := data["data"].(map[string]interface{})
	assert.Equal(t, "Andrew", objData["firstName"])
	assert.Equal(t, "Wansley", objData["lastName"])
	assert.Equal(t, 32.0, objData["age"])
}
