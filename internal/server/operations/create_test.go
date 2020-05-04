package operations

import (
	"awans.org/aft/internal/server/db"
	"awans.org/aft/internal/server/middleware"
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

func TestCreateServerParseSimple(t *testing.T) {
	appDB := db.New()
	db.AddSampleModels(appDB)
	assert := assert.New(t)
	req, err := http.NewRequest("POST", "/user.create", strings.NewReader(
		`{"data":{
			"firstName":"Andrew",
			"lastName":"Wansley",
			"age": 32,
		}}`))
	if err != nil {
		t.Fatal(err)
	}
	req = mux.SetURLVars(req, map[string]string{"object": "user"})

	qs := CreateServer{}
	ctx := middleware.NewRWTxContext(context.Background(), appDB.NewRWTx())
	parsed, err := qs.Parse(ctx, req)
	if err != nil {
		t.Error(err)
	}
	parsedReq := parsed.(CreateRequest)

	assert.IsType(parsedReq, CreateRequest{})

	op := parsedReq.Operation
	st := op.Struct

	reader := dynamicstruct.NewReader(st)
	u := reader.GetField("Id").Interface().(uuid.UUID)

	assert.Zero(u)

	firstName := reader.GetField("Firstname").String()
	assert.Equal(firstName, "Andrew")

	age := reader.GetField("Age").Int()
	assert.Equal(age, 32)
}

func TestCreateServerServe(t *testing.T) {
	appDB := db.New()
	db.AddSampleModels(appDB)
	jsonString := `{ "firstName":"Andrew", 
"lastName":"Wansley",
"age": 32}`
	u := makeStruct(appDB.NewTx(), "user", jsonString)
	cOp := db.CreateOperation{
		Struct: u,
		Nested: []db.NestedOperation{},
	}
	req := CreateRequest{Operation: cOp}
	cs := CreateServer{}
	ctx := middleware.NewRWTxContext(context.Background(), appDB.NewRWTx())
	resp, err := cs.Serve(ctx, req)
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
