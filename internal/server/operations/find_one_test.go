package operations

import (
	"awans.org/aft/internal/server/db"
	"context"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/ompluscator/dynamic-struct"
	"github.com/stretchr/testify/assert"
	"net/http"
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
	db.AddSampleModels(appDB)
	req, err := http.NewRequest("POST", "/user.query", strings.NewReader(
		`{"where": {
		"id": "2b1e9f08-38a9-4a36-b653-0fa0cbc8cad2"
		}}`))
	if err != nil {
		t.Fatal(err)
	}
	req = mux.SetURLVars(req, map[string]string{"object": "user"})
	fos := FindOneServer{DB: appDB}
	ifc, err := fos.Parse(context.Background(), req)
	_, ok := ifc.(FindOneRequest)
	if !ok {
		t.Fatal("Didn't return a FindOneRequest")
	}
}

func TestFindOneServerServe(t *testing.T) {
	appDB := db.New()
	db.AddSampleModels(appDB)
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
	resp, err := fos.Serve(context.Background(), req)
	if err != nil {
		t.Error(err)
	}
	bytes, _ := json.Marshal(resp)

	var data map[string]interface{}
	json.Unmarshal(bytes, &data)
	objData := data["data"].(map[string]interface{})
	assert.Equal(t, "Andrew", objData["firstName"])
	assert.Equal(t, "Wansley", objData["lastName"])
	assert.Equal(t, 32.0, objData["age"])
}
