package starlark

import (
	"awans.org/aft/internal/bizdatatypes"
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
	{`y = FindOne("code", Eq("name", "emailAddressValidator"))
result(Exec(y, "chase@hensel.com"))`, "chase@hensel.com", false},
	{`y = FindOne("code", Eq("name", "emailAddressValidator"))
result(Exec(y, "sue@"))`, "", true},
}

func TestDB(t *testing.T) {
	appDB := db.NewTest()
	tx := appDB.NewRWTx()
	r1 := db.RecordForModel(db.CodeModel)
	db.SaveCode(r1, bizdatatypes.EmailAddressValidator)
	tx.Insert(r1)
	tx.Commit()
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
