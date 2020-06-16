package runtime

import (
	"awans.org/aft/internal/api"
	"awans.org/aft/internal/bus"
	"awans.org/aft/internal/db"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

var Andrew = db.Datatype{
	ID:        uuid.MustParse("46c0ee11-3943-452d-9420-925dd9be8208"),
	Name:      "andrew",
	StoredAs:  db.StringStorage,
	Validator: AndrewCode,
}

var testingRox = db.Datatype{
	ID:        uuid.MustParse("9d792f82-018e-47d1-a2e5-a1b5b4822fd9"),
	Name:      "testingRox",
	StoredAs:  db.StringStorage,
	Validator: testingRoxCode,
}

var AndrewCode = db.Code{
	ID:                uuid.MustParse("a4615a60-afed-4f29-b674-e24f35618847"),
	Name:              "andrew",
	Runtime:           db.Starlark,
	FunctionSignature: db.FromJSON,
	Code: `def validator(arg):
	if str(arg) == "Andrew":
  		return "testing rox"
	error("arg should be Andrew!!!")
`,
}

var testingRoxCode = db.Code{
	ID:                uuid.MustParse("5b0cfd40-4f3d-4890-b3a9-923ab8740043"),
	Name:              "testingRox",
	Runtime:           db.Starlark,
	FunctionSignature: db.FromJSON,
	Code: `def validator(arg):
	return "testing rox"`,
}

var UserStarlark = db.Model{
	ID:   uuid.MustParse("c1da149d-8ba0-429a-ab66-a8f2973c9e1e"),
	Name: "starlark",
	Attributes: map[string]db.Attribute{
		"firstName": db.Attribute{
			ID:       uuid.MustParse("a6e4e877-3e8c-4b17-9e80-0b896c0a1086"),
			Datatype: Andrew,
		},
		"lastName": db.Attribute{
			ID:       uuid.MustParse("9bd6a83c-b805-4daf-b56f-6824f51fdbca"),
			Datatype: testingRox,
		},
	},
}

var starlarkTests = []struct {
	in          string
	out         string
	field       string
	shouldError bool
}{
	{`{"data":{
			"firstName":"Chase",
			"lastName":"Wansley"
		}
	}`, "", "", true},
	{`{"data":{
			"firstName":"Andrew",
			"lastName":"Wansley"
		}
	}`, "testing rox", "lastName", false},
}

func TestError(t *testing.T) {
	tt := starlarkTests[0]
	runner(t, tt.in, tt.out, tt.field, tt.shouldError)
}

func TestNoError(t *testing.T) {
	tt := starlarkTests[1]
	runner(t, tt.in, tt.out, tt.field, tt.shouldError)
}

func runner(t *testing.T, in, out, field string, shouldError bool) {
	appDB := db.New(&Executor{})
	tx := appDB.NewRWTx()

	r := db.RecordForModel(db.CodeModel)
	db.SaveCode(r, AndrewCode)
	tx.Insert(r)

	r = db.RecordForModel(db.CodeModel)
	db.SaveCode(r, testingRoxCode)
	tx.Insert(r)

	r = db.RecordForModel(db.DatatypeModel)
	db.SaveDatatype(r, Andrew)
	tx.Insert(r)

	r = db.RecordForModel(db.DatatypeModel)
	db.SaveDatatype(r, testingRox)
	tx.Insert(r)

	tx.SaveModel(UserStarlark)
	tx.Commit()

	eventbus := bus.New()
	db.AddSampleModels(appDB)
	req, err := http.NewRequest("POST", "/starlark.create", strings.NewReader(in))
	req = mux.SetURLVars(req, map[string]string{"modelName": "starlark"})
	cs := api.CreateHandler{DB: appDB, Bus: eventbus}
	w := httptest.NewRecorder()
	err = cs.ServeHTTP(w, req)
	if shouldError {
		assert.Error(t, err)
		return
	}
	var data map[string]interface{}
	result := w.Result()
	bytes, err := ioutil.ReadAll(result.Body)
	result.Body.Close()
	if err != nil {
		t.Error(err)
	}
	json.Unmarshal(bytes, &data)
	objData := data["data"].(map[string]interface{})
	assert.Equal(t, out, objData[field])
}
