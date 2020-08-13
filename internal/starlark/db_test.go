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
	{`x = aft.api.create("code", {"data" : {"name" : "bob"}})
y = aft.api.findOne("code", {"where" : {"name": "bob"}})
def main():
    return y.name`, "bob", false},
	{`y = aft.api.delete("code", {"where" : {"name": "bob"}})
z = aft.api.findOne("code", {"where" : {"name": "bob"}})
def main():
    return z.name`, "", true},
	{`aft.api.create("code", {"data" : {"name" : "bob"}})
aft.api.update("code", {"where" : {"name" : "bob"}, "data" : {"name" : "sue"}})
z = aft.api.findOne("code", {"where" : {"name": "sue"}})
def main():
    return z.name`, "sue", false},
	{`def main():
   return str(aft.function.int(5))`, "5", false},
	{`def main():
    return aft.function.int("sue@")`, "", true},
}

func TestDB(t *testing.T) {
	appDB := db.NewTest()
	tx := appDB.NewRWTx()
	for _, tt := range dbTests {
		fh := MakeStarlarkFunction(db.NewID(), "", db.RPC, tt.in)
		r, err := fh.CallWithEnv("", DBLib(tx))
		if tt.shouldError {
			assert.Error(t, err)
		} else {
			assert.Equal(t, tt.out, r)
		}
	}
}
