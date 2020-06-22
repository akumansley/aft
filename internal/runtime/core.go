package runtime

import (
	"awans.org/aft/internal/datatypes"
	"awans.org/aft/internal/db"
	"awans.org/aft/internal/starlark"
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
		fs := c.FunctionSignature
		fh = &starlark.StarlarkFunctionHandle{Code: code, FunctionSignature: fs}
	}
	return fh
}

type Executor struct{}

func (*Executor) Invoke(c db.Code, args interface{}) (interface{}, error) {
	return codeToFunctionHandle(c).Invoke(args)
}
