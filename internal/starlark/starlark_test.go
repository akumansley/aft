package starlark

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

var runtimeTests = []struct {
	in          string
	out         interface{}
	shouldError bool
}{
	{`result(5+5)`, int64(10), false},
	{"result('bob' + '_burgers')", "bob_burgers", false},
	{"result(4.7 + 5.1)", 9.8, false},
}

func TestStarlark(t *testing.T) {
	for _, tt := range runtimeTests {
		fh := StarlarkFunctionHandle{Code: tt.in}
		r, err := fh.Invoke("")
		if tt.shouldError {
			assert.Error(t, err)
		} else {
			assert.Equal(t, tt.out, r)
		}
	}
}
