package lib

import (
	"go.starlark.net/starlark"
	"go.starlark.net/starlarkstruct"
)

var Lib = starlark.StringDict{
	"re":           re,
	"urlparse":     urlparse,
	"struct":       starlark.NewBuiltin("struct", starlarkstruct.Make),
	"error":        starlark.NewBuiltin("error", makeError),
	"loadFunction": loadFunction,
}
