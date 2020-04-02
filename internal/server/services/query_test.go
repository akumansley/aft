package services

import (
	"awans.org/aft/internal/server/db"
	"github.com/gorilla/mux"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestQueryServerParse(t *testing.T) {
	db.SetupTestData()
	req, err := http.NewRequest("POST", "/objects.query", strings.NewReader(
		`{
		"query": "Cekw67uyMpBGZLRP2HFVbe"
		}`))
	if err != nil {
		t.Fatal(err)
	}
	req = mux.SetURLVars(req, map[string]string{"object": "objects"})
	qs := QueryServer{}
	parsedReq, ok := qs.Parse(req).(QueryRequest)
	if !ok {
		t.Fatal("Didn't return a QueryRequest")
	}
	if parsedReq.Type != "objects" {
		t.Errorf("Expected a type of objects, got %v", parsedReq.Type)
	}
	if parsedReq.Q != "Cekw67uyMpBGZLRP2HFVbe" {
		t.Errorf("Didn't parse query as expected; got %v", parsedReq.Q)
	}
}

func TestQueryServerServe(t *testing.T) {
	db.SetupTestData()
	req := QueryRequest{Q: "Cekw67uyMpBGZLRP2HFVbe", Type: "objects"}
	qs := QueryServer{}
	rr := httptest.NewRecorder()
	qs.Serve(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
	expected := `{"data":[{"id":"Cekw67uyMpBGZLRP2HFVbe","name":"Test","fields":[{"name":"f1","type":2}]}]}`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}
