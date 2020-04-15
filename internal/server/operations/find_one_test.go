package operations

import (
	"awans.org/aft/internal/server/db"
	"github.com/gorilla/mux"
	"net/http"
	"strings"
	"testing"
)

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

func TestQueryServerServe(t *testing.T) {
	// TODO come back to this after we get writes up
	// db.SetupSchema()
	// req := QueryRequest{Q: "Cekw67uyMpBGZLRP2HFVbe", Type: "objects"}
	// qs := QueryServer{}
	// rr := httptest.NewRecorder()
	// qs.Serve(rr, req)
	// if status := rr.Code; status != http.StatusOK {
	// 	t.Errorf("handler returned wrong status code: got %v want %v",
	// 		status, http.StatusOK)
	// }
	// expected := `{"data":[{"id":"Cekw67uyMpBGZLRP2HFVbe","name":"Test","fields":[{"name":"f1","type":2}]}]}`
	// if rr.Body.String() != expected {
	// 	t.Errorf("handler returned unexpected body: got %v want %v",
	// 		rr.Body.String(), expected)
	// }
}
