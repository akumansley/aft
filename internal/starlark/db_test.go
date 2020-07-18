package starlark

import (
	"awans.org/aft/internal/db"
	"github.com/stretchr/testify/assert"
	"testing"
)

var dbTests = []struct {
	in          string
	out         interface{}
	shouldError bool
}{
	{`x = Insert("code", {"name" : "bob"})
y = FindOne("code", Eq("name", "bob"))
result(y.Get("name"))`, "bob", false},
	{`y = FindOne("code", Eq("name", "bob"))
Delete(y)
z = FindOne("code", Eq("name", "bob"))
result(z.Get("name"))`, "", true},
	{`x = Insert("code", {"name" : "bob"})
Update(x, {"name": "sue"})
z = FindOne("code", Eq("name", "sue"))
result(z.Get("name"))`, "sue", false},
	{`y = FindOne("code", Eq("name", "int"))
out, ran = Exec(y, "5")
result(out)`, "5", false},
	{`y = FindOne("code", Eq("name", "int"))
out, ran = Exec(y, "sue@")
result(out)`, "invalid value for type: expected int got string", false},
}

func TestDB(t *testing.T) {
	appDB := db.NewTest()
	tx := appDB.NewRWTx()
	for _, tt := range dbTests {
		fh := StarlarkFunctionHandle{Code: tt.in, Env: DBLib(tx)}
		r, err := fh.Invoke("")
		if tt.shouldError {
			assert.Error(t, err)
		} else {
			assert.Equal(t, tt.out, r)
		}
	}
}
