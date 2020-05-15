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
	fos := FindOneServer{}
	ctx := middleware.NewTxContext(context.Background(), appDB.NewTx())
	ifc, err := fos.Parse(ctx, req)
	_, ok := ifc.(FindOneRequest)
	if !ok {
		t.Fatal("Didn't return a FindOneRequest")
	}
}

func TestFindOneServerServe(t *testing.T) {
	appDB := db.New()
	db.AddSampleModels(appDB)
	jsonString := `{ "firstName":"Andrew", "lastName":"Wansley", "age": 32}`
	u := makeRecord(appDB.NewTx(), "user", jsonString)
	cOp := CreateOperation{
		Record: u,
		Nested: []NestedOperation{},
	}
	tx := appDB.NewRWTx()
	cOp.Apply(tx)
	tx.Commit()

	req := FindOneRequest{Operation: FindOneOperation{
		ModelName: "user",
		UniqueQuery: UniqueQuery{
			Key: "Id",
			Val: u.Id(),
		},
	}}
	fos := FindOneServer{}
	ctx := middleware.NewTxContext(context.Background(), appDB.NewTx())
	resp, err := fos.Serve(ctx, req)
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
