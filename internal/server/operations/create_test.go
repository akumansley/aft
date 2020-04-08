package operations

import (
	"awans.org/aft/internal/server/db"
	"github.com/gorilla/mux"
	"net/http"
	// "net/http/httptest"
	"github.com/google/uuid"
	"github.com/ompluscator/dynamic-struct"
	"github.com/stretchr/testify/assert"
	"os"
	"strings"
	"testing"
)

func TestMain(m *testing.M) {
	db.InitDB()
	AddSampleModels()
	os.Exit(m.Run())
}

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

func TestCreateServerParseNestedCreate(t *testing.T) {
	assert := assert.New(t)
	req, err := http.NewRequest("POST", "/user.create", strings.NewReader(
		`{"data":{
			"id":"15852d31-3bd4-4fc4-abd0-e4c7497644ab",
			"firstName":"Andrew",
			"lastName":"Wansley",
			"age": 32,
			"profile": {
			  "create": {
			    "id": "c8f857ca-204c-46ab-a96e-d69c1df2fa4f",
			    "text": "My bio.."
			  }
			}
		}
	}`))
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

	reader := dynamicstruct.NewReader(st)
	u := reader.GetField("Id").Interface().(uuid.UUID)
	assert.Equal(u.String(), "15852d31-3bd4-4fc4-abd0-e4c7497644ab")

	firstName := reader.GetField("Firstname").String()
	assert.Equal(firstName, "Andrew")

	age := reader.GetField("Age").Int()
	assert.Equal(age, 32)

	nested := op.Nested
	assert.Len(nested, 1)

	profileCreate := nested[0].(NestedCreateOperation)
	parsedProfile := profileCreate.Struct

	reader = dynamicstruct.NewReader(parsedProfile)

	u = reader.GetField("Id").Interface().(uuid.UUID)
	assert.Equal(u.String(), "c8f857ca-204c-46ab-a96e-d69c1df2fa4f")

	bio := reader.GetField("Text").String()
	assert.Equal(bio, "My bio..")
}

func TestCreateServerParseNestedCreateMany(t *testing.T) {
	assert := assert.New(t)
	req, err := http.NewRequest("POST", "/user.create", strings.NewReader(
		`{"data":{
			"id":"15852d31-3bd4-4fc4-abd0-e4c7497644ab",
			"firstName":"Andrew",
			"lastName":"Wansley",
			"age": 32,
			"posts": {
			  "create": [{
			    "id": "57e3f538-d35a-45e8-acdf-0ab916d8194f",
			    "text": "post1"
			  }, {
			    "id": "6327fe0e-c936-4332-85cd-f1b42f6f337a",
			    "text": "post2"
			  }]
			}
		}
	}`))
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

	reader := dynamicstruct.NewReader(st)
	u := reader.GetField("Id").Interface().(uuid.UUID)
	assert.Equal(u.String(), "15852d31-3bd4-4fc4-abd0-e4c7497644ab")
	firstName := reader.GetField("Firstname").String()
	assert.Equal(firstName, "Andrew")
	age := reader.GetField("Age").Int()
	assert.Equal(age, 32)
	nested := op.Nested
	assert.Len(nested, 2)

	postCreate := nested[0].(NestedCreateOperation)
	parsedPost := postCreate.Struct

	reader = dynamicstruct.NewReader(parsedPost)
	u = reader.GetField("Id").Interface().(uuid.UUID)
	assert.Equal(u.String(), "57e3f538-d35a-45e8-acdf-0ab916d8194f")
	bio := reader.GetField("Text").String()
	assert.Equal(bio, "post1")

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
