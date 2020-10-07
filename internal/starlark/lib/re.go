package lib

import (
	"fmt"
	"regexp"

	"go.starlark.net/starlark"
	"go.starlark.net/starlarkstruct"
)

var re = &starlarkstruct.Module{
	Name: "re",
	Members: starlark.StringDict{
		"compile": starlark.NewBuiltin("compile", compile),
	},
}

func compile(thread *starlark.Thread, b *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (out starlark.Value, err error) {
	var val starlark.Value
	if err := starlark.UnpackPositionalArgs(b.Name(), args, kwargs, 1, &val); err != nil {
		return starlark.None, err
	}
	patternS, ok := starlark.AsString(val)
	if !ok {
		return starlark.None, fmt.Errorf("Invalid string: %s", val)
	}

	//here is the main point of the function
	re, err := regexp.Compile(patternS)
	if err != nil {
		return starlark.None, err
	}

	var match = func(thread *starlark.Thread, b *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (out starlark.Value, err error) {
		var val starlark.Value
		if err := starlark.UnpackPositionalArgs(b.Name(), args, kwargs, 1, &val); err != nil {
			return out, err
		}
		s, ok := starlark.AsString(val)
		if !ok {
			return starlark.None, fmt.Errorf("Invalid string: %s", val)
		}
		//this is where we match against a string
		ok = re.MatchString(s)
		if ok {
			return starlark.True, nil
		}
		return starlark.False, nil
	}
	return &starlarkstruct.Module{
		Name: b.Name(),
		Members: starlark.StringDict{
			"match": starlark.NewBuiltin("match", match),
		},
	}, nil
}
