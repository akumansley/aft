package operations

import (
	"github.com/gorilla/mux"
	"net/http"
	// "net/http/httptest"
	"github.com/google/uuid"
	"github.com/ompluscator/dynamic-struct"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestCreateServerParseSimple(t *testing.T) {
	assert := assert.New(t)
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
