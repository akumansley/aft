package runtime

import (
	"awans.org/aft/internal/datatypes"
	"awans.org/aft/internal/db"
	"awans.org/aft/internal/starlark"
	"fmt"
)

type FunctionHandle interface {
	Invoke(arg interface{}) (interface{}, error)
}

func codeToFunctionHandle(c db.Code) FunctionHandle {
	var fh FunctionHandle
	switch c.Runtime {
	case db.Golang:
		fh = &datatypes.GoFunctionHandle{Function: c.Function}
	case db.Starlark:
		code := c.Code
		if c.FunctionSignature == db.FromJSON {
			code = fmt.Sprintf("%s\n\nresult(validator(args))", code)
		}
		fh = &starlark.StarlarkFunctionHandle{Code: code}
	}
	return fh
}

type Executor struct{}

func (*Executor) Invoke(c db.Code, args interface{}) (interface{}, error) {
	return codeToFunctionHandle(c).Invoke(args)
}
