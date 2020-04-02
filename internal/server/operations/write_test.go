package operations

import (
	"awans.org/aft/internal/server/db"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestWriteServerParse(t *testing.T) {
	db.SetupTestData()
	req, err := http.NewRequest("POST", "/objects.create", strings.NewReader(
		`{
		"body":{
			"id":"abc123",
			"name":"Test",
			"fields":[{"name":"f1","type":2}]}
		}`))
	if err != nil {
		t.Fatal(err)
	}
	req = mux.SetURLVars(req, map[string]string{"object": "objects"})
	qs := WriteServer{}
	parsedReq, ok := qs.Parse(req).(WriteRequest)
	if !ok {
		t.Fatal("Didn't return a WriteRequest")
	}
	if parsedReq.Type != "objects" {
		t.Errorf("Expected a type of objects, got %v", parsedReq.Type)
	}
	if parsedReq.Body["id"].(string) != "abc123" {
		t.Errorf("Didn't parse id as expected; got %v", parsedReq.Body["id"])
	}
}

func TestWriteServerServe(t *testing.T) {
	db.SetupTestData()
	data := map[string]interface{}{
		"id":   "abc123",
		"name": "Test3",
		"fields": []interface{}{
			map[string]interface{}{
				"name": "f1",
				"type": 2,
			},
		},
	}
	req := WriteRequest{Body: data, Type: "objects"}
	cs := WriteServer{}
	rr := httptest.NewRecorder()
	cs.Serve(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
	expected := `{"body":{"id":"abc123","name":"Test3","fields":[{"name":"f1","type":2}]}}`
	assert.JSONEq(t, expected, rr.Body.String())
}