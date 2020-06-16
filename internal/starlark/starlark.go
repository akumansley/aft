package starlark

import (
	"fmt"
	"go.starlark.net/starlark"
)

type StarlarkFunctionHandle struct {
	Code   string
	Env    map[string]interface{}
	result interface{}
	err    interface{}
}

func (s *StarlarkFunctionHandle) Invoke(input interface{}) (interface{}, error) {
	globals, err := s.createGlobals(input)
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
	return s.result, nil
}