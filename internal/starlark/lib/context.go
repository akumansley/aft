package lib

import (
	"context"
	"errors"

	"go.starlark.net/starlark"
)

type ContextValue struct {
	context.Context
}

func (c ContextValue) String() string {
	return "context()"
}

func (c ContextValue) Type() string {
	return "context.Context"
}

func (c ContextValue) Freeze() {}

func (c ContextValue) Truth() starlark.Bool {
	return starlark.Bool(c.Context != nil)
}

func (c ContextValue) Hash() (uint32, error) {
	return 0, errors.New("Unhashable")
}
