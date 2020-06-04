package api

import (
	"awans.org/aft/internal/bus"
	"awans.org/aft/internal/db"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

var Andrew = db.Datatype{
	ID:          uuid.MustParse("46c0ee11-3943-452d-9420-925dd9be8208"),
	Name:        "andrew",
	Validator:   AndrewValidator,
	StorageType: db.StringType,
}

var AndrewValidator = db.Code{
	ID:       uuid.MustParse("aaea187b-d153-4c4a-a7e7-cda599b02ba6"),
	Name:     "andrewValidator",
	Runtime:  db.Starlark,
	Function: db.Validator,
	Code: `
def func(arg):
  if arg == "Andrew":
  	return ""
  return "arg should be Andrew!!!"
arg.Error = func(arg.Value)
`,
}

var UserStarlark = db.Model{
	ID:   uuid.MustParse("887a91b8-3857-4b4d-a633-a6386a4fae25"),
	Name: "userStarlark",
	Attributes: map[string]db.Attribute{
		"firstName": db.Attribute{
			ID:       uuid.MustParse("2afdc6d7-9715-41eb-80d0-20b5132efe94"),
			Datatype: Andrew,
		},
		"lastName": db.Attribute{
			ID:       uuid.MustParse("462212e7-dd94-403e-8314-e271fd7ccec9"),
			Datatype: db.String,
		},
		"age": db.Attribute{
			ID:       uuid.MustParse("7b0f19ab-a615-49f7-b5a6-d2054d442a76"),
			Datatype: db.Int,
		},
		"emailAddress": db.Attribute{
			ID:       uuid.MustParse("0fe6bd01-9828-43ac-b004-37620083344d"),
			Datatype: db.EmailAddress,
		},
	},
	LeftRelationships: []db.Relationship{
		db.UserPosts,
		db.UserProfile,
	},
}

func TestCreateAndrewType(t *testing.T) {
	appDB := db.New()
	tx := appDB.NewRWTx()
	r := db.RecordForModel(db.DatatypeModel)
	db.SaveDatatype(r, Andrew)
	tx.Insert(r)

	r = db.RecordForModel(db.CodeModel)
	db.SaveCode(r, AndrewValidator)
	tx.Insert(r)
	tx.SaveModel(UserStarlark)
	tx.Commit()

	eventbus := bus.New()
	db.AddSampleModels(appDB)
	req, err := http.NewRequest("POST", "/userStarlark.create", strings.NewReader(
		`{"data":{
			"firstName":"Chase",
			"lastName":"Wansley",
			"age": 32,
			"emailAddress": "andrew.wansley@gmail.com",
			"profile": {
				"create": {
					"text": "hello"
				}
			}
		},
		"include": {
			"profile": true
		}

	}`))

	req = mux.SetURLVars(req, map[string]string{"modelName": "userStarlark"})

	cs := CreateHandler{db: appDB, bus: eventbus}
	w := httptest.NewRecorder()
	err = cs.ServeHTTP(w, req)
	assert.Error(t, err)
}
