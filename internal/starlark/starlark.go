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
	msgs              string
}

func (s *StarlarkFunctionHandle) Invoke(input interface{}) (interface{}, error) {
	inp, err := convert.ToValue(input)
	if err != nil {
		return nil, err
	}
	globals, err := s.createEnv()
	if err != nil {
		return nil, err
	}
	// Run the starlark interpreter!
	th := &starlark.Thread{Load: nil}
	globals, err = starlark.ExecFile(th, "", []byte(s.Code), globals)
	if err != nil {
		return nil, err
	}
	if globals["main"] == nil {
		return nil, fmt.Errorf("Missing main function")
	}
	// Check how many args main takes
	numArgs := (globals["main"].(*starlark.Function)).NumParams()
	if numArgs > 1 {
		return nil, fmt.Errorf("Main can't take more than 1 arg")
	}
	var args []starlark.Value
	if numArgs == 1 {
		args = append(args, inp)
	}
	out, err := starlark.Call(th, globals["main"], args, nil)
	if err != nil {
		return nil, err
	}
	// If there were print statements, print them
	if s.msgs != "" {
		return fmt.Sprintf("%s%v", s.msgs, recursiveFromValue(out)), nil
	}
	return recursiveFromValue(out), nil
}
