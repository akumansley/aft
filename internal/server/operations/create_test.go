package operations

import (
	"awans.org/aft/internal/server/db"
	"github.com/gorilla/mux"
	"net/http"
	// "net/http/httptest"
	"fmt"
	"github.com/google/uuid"
	"github.com/ompluscator/dynamic-struct"
	"strings"
	"testing"
)

func TestCreateServerParseSimple(t *testing.T) {
	db.InitDB()
	AddSampleModels()

	req, err := http.NewRequest("POST", "/user.create", strings.NewReader(
		`{"data":{
			"id":"15852d31-3bd4-4fc4-abd0-e4c7497644ab",
			"firstName":"Andrew",
			"lastName":"Wansley",
			"age": 32,
		}}`))
	if err != nil {
		t.Fatal(err)
	}
	req = mux.SetURLVars(req, map[string]string{"object": "user"})
	qs := CreateServer{}
	parsedReq, ok := qs.Parse(req).(CreateRequest)
	if !ok {
		t.Fatal("Didn't return a CreateRequest")
	}
	op := parsedReq.Operation
	st := op.Struct
	fmt.Printf("struct is %v", st)
	reader := dynamicstruct.NewReader(st)
	uuid := reader.GetField("Id").Interface().(uuid.UUID)

	if uuid.String() != "15852d31-3bd4-4fc4-abd0-e4c7497644ab" {
		t.Errorf("Didn't parse id as expected; got %v", uuid)
	}
}

func TestCreateServerParseNestedCreate(t *testing.T) {
	// db.InitDB()
	// AddSampleModels()

	// req, err := http.NewRequest("POST", "/person.create", strings.NewReader(
	// 	`{"data":{
	// 		"id":"15852d31-3bd4-4fc4-abd0-e4c7497644ab",
	// 		"firstName":"Andrew",
	// 		"lastName":"Wansley",
	// 		"age": 32,
	// 		"profile": {
	// 		  "create": {
	// 		    "id": "c8f857ca-204c-46ab-a96e-d69c1df2fa4f",
	// 		    "text": "My bio.."
	// 		  }
	// 		}
	// 	}
	// 	}`))
	// if err != nil {
	// 	t.Fatal(err)
	// }
	// req = mux.SetURLVars(req, map[string]string{"object": "person"})
}

// func TestWriteServerServe(t *testing.T) {
// 	db.SetupSchema()
// 	data := map[string]interface{}{
// 		"id":   "abc123",
// 		"name": "Test3",
// 		"fields": []interface{}{
// 			map[string]interface{}{
// 				"name": "f1",
// 				"type": 2,
// 			},
// 		},
// 	}
// 	req := WriteRequest{Body: data, Type: "objects"}
// 	cs := WriteServer{}
// 	rr := httptest.NewRecorder()
// 	cs.Serve(rr, req)
// 	if status := rr.Code; status != http.StatusOK {
// 		t.Errorf("handler returned wrong status code: got %v want %v",
// 			status, http.StatusOK)
// 	}
// 	expected := `{"body":{"id":"abc123","name":"Test3","fields":[{"name":"f1","type":2}]}}`
// 	assert.JSONEq(t, expected, rr.Body.String())
// }
