package lib

import (
	"net/url"

	"go.starlark.net/starlark"
)

var urlparse = starlark.NewBuiltin("urlparse", urlParse)

func urlParse(thread *starlark.Thread, b *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	var s string

	if err := starlark.UnpackArgs(b.Name(), args, kwargs, "s", &s); err != nil {
		return nil, err
	}

	out, err := url.ParseRequestURI(s)

	if err != nil {
		return starlark.String(""), nil
	}

	return starlark.String(out.String()), nil
}
