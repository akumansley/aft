package starlark

import (
	"awans.org/aft/internal/db"
	"fmt"
	"go.starlark.net/starlark"
)

type StarlarkFunctionHandle struct {
	Code              string
	FunctionSignature db.FunctionSignature
	Env               map[string]interface{}
	result            interface{}
	err               interface{}
}

func (s *StarlarkFunctionHandle) Invoke(input interface{}) (interface{}, error) {
	i, err := recursiveToValue(input)
	if err != nil {
		return nil, err
	}
	globals, err := s.createEnv(i)
	if err != nil {
		return nil, err
	}

	if s.FunctionSignature == db.FromJSON {
		s.Code = fmt.Sprintf("%s\n\nresult(validator(args))", s.Code)
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
