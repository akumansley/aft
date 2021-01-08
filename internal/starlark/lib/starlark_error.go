package lib

import (
	"errors"
	"fmt"

	"go.starlark.net/starlark"
)

type StarlarkError struct {
	Code    string
	Message string
}

func (s *StarlarkError) String() string {
	return fmt.Sprintf("error(code=\"%v\", message=\"%v\")", s.Code, s.Message)
}

func (s *StarlarkError) Type() string {
	return "StarlarkError"
}

func (s *StarlarkError) Freeze() {}

func (s *StarlarkError) Truth() starlark.Bool {
	return starlark.Bool(false)
}

func (s *StarlarkError) Hash() (uint32, error) {
	return 0, errors.New("Unhashable")
}

func makeError(_ *starlark.Thread, _ *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	var code string
	var message string

	if len(args) > 0 {
		return nil, fmt.Errorf("struct: unexpected positional arguments")
	}
	err := starlark.UnpackArgs("error", args, kwargs, "code", &code, "message", &message)
	return &StarlarkError{Code: code, Message: message}, err
}
