package repl

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

var runtimeTests = []struct {
	in  string
	out interface{}
}{
	{"result(5+5)", "10"},
	{"result('bob' + '_burgers')", "bob_burgers"},
	{"result(4.7 + 5.1)", "9.8"},
	{"x = 5", ""},
}

func TestStarlark(t *testing.T) {
	for _, tt := range runtimeTests {
		r := eval(tt.in, nil)
		assert.Equal(t, tt.out, r)
	}
}
