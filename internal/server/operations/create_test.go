package operations

import (
	"awans.org/aft/internal/db"
	"awans.org/aft/internal/server/middleware"
	"context"
	"encoding/json"
	"github.com/gorilla/mux"
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
	rec := op.Record

	u := rec.Id()

	assert.Zero(u)

	firstName := rec.Get("firstName").(string)
	assert.Equal(firstName, "Andrew")

	age := rec.Get("Age").(int64)
	assert.Equal(age, int64(32))
}

func TestCreateServerServe(t *testing.T) {
	appDB := db.New()
	db.AddSampleModels(appDB)
	jsonString := `{ "firstName":"Andrew", "lastName":"Wansley", "age": 32}`
	u := makeRecord(appDB.NewTx(), "user", jsonString)
	cOp := CreateOperation{
		Record: u,
		Nested: []NestedOperation{},
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
