package bizdatatypes

import (
	"awans.org/aft/internal/db"
	"awans.org/aft/internal/starlark"
	"github.com/stretchr/testify/assert"
	"testing"
)

var runtimeTests = []struct {
	in          string
	out         interface{}
	shouldError bool
}{
	{`def main():
   return 5+5`, int64(10), false},
	{`def main():
   return 'bob' + '_burgers'`, "bob_burgers", false},
	{`def main():
   return 4.7 + 5.1`, 9.8, false},
}

func TestStarlark(t *testing.T) {
	for _, tt := range runtimeTests {
		sf := starlark.MakeStarlarkFunction(db.NewID(), "test", db.RPC, tt.in)
		r, err := sf.Call("")
		if tt.shouldError {
			assert.Error(t, err)
		} else {
			assert.Equal(t, tt.out, r)
		}
	}
}