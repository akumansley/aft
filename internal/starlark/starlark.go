package starlark

import (
	"awans.org/aft/internal/db"
	"fmt"
	"github.com/chasehensel/starlight/convert"
	"go.starlark.net/resolve"
	"go.starlark.net/starlark"
)

//configure starlark
func init() {
	resolve.AllowNestedDef = true // allow def statements within function bodies
	resolve.AllowLambda = true    // allow lambda expressions
	resolve.AllowFloat = true     // allow floating point literals, the 'float' built-in, and x / y
	resolve.AllowSet = true       // allow the 'set' built-in
	resolve.AllowRecursion = true // allow while statements and recursive functions
}

type StarlarkFunctionHandle struct {
	Code              string
	FunctionSignature db.FunctionSignatureEnumValue
	Env               map[string]interface{}
	result            interface{}
	err               interface{}
}

func (s *StarlarkFunctionHandle) Invoke(input interface{}) (interface{}, error) {
	i, err := convert.ToValue(input)
	if err != nil {
		return nil, err
	}
	globals, err := s.createEnv(i)
	if err != nil {
		return nil, err
	}

	// Run the starlark interpreter!
	_, err = starlark.ExecFile(&starlark.Thread{Load: nil}, "", []byte(s.Code), globals)
	if err != nil {
		if evale, ok := err.(*starlark.EvalError); ok {
			return nil, fmt.Errorf("\n%s", evale.Backtrace())
		}
		return nil, err
	}
	if s.err != nil {
		return nil, fmt.Errorf("Raised: %s", s.err)
	}
	out := recursiveFromValue(s.result)
	return out, nil
}
