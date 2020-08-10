package starlark

import (
	"fmt"
	"github.com/chasehensel/starlight/convert"
	"go.starlark.net/starlark"
	"go.starlark.net/starlarkstruct"
	"net/url"
	"reflect"
)

var (
	ErrInvalidInput = fmt.Errorf("Bad input:")
)

func assertString(val interface{}) (string, error) {
	if reflect.TypeOf(val) != reflect.TypeOf("") {
		return "", fmt.Errorf(" %w expected string, but found %T", ErrInvalidInput, val)
	}
	return val.(string), nil
}

func StdLib(c *call) map[string]interface{} {
	env := map[string]interface{}{
		"re": &starlarkstruct.Module{
			Name: "re",
			Members: starlark.StringDict{
				"compile": starlark.NewBuiltin("compile", compile),
			},
		},
		"test": func(str interface{}, a ...interface{}) {
			input := fmt.Sprintf("%v", str)
			fmt.Printf(input, a...)
		},
		"print": func(a ...interface{}) {
			parser := ""
			for _, _ = range a {
				parser = fmt.Sprintf("%s%s", parser, "%v")
			}
			c.msgs = fmt.Sprintf("%s%s\n", c.msgs, fmt.Sprintf(parser, a...))
		},
		"urlparse": func(u interface{}) (string, bool) {
			us, err := assertString(u)
			if err != nil {
				return "", false
			}
			out, err := url.ParseRequestURI(us)
			if err != nil {
				return "", false
			}
			return fmt.Sprintf("%s", out), true
		},
		"sprint": func(str string, a ...interface{}) string { return fmt.Sprintf(str, a...) },
	}
	return env
}

type call struct {
	Env  map[string]interface{}
	msgs string
}

func CreateEnv(c *call) (starlark.StringDict, error) {
	stdlib := StdLib(c)
	env, err := convert.MakeStringDict(stdlib)
	if err != nil {
		return nil, err
	}
	api, err := convert.MakeStringDict(c.Env)
	if err != nil {
		return nil, err
	}
	//API overwrites local
	env = clobber(env, api)
	env = union(env, api)

	return env, nil
}

func union(a, b starlark.StringDict) starlark.StringDict {
	if a == nil {
		return b
	}
	for k, v := range b {
		a[k] = v
	}
	return a
}

func clobber(a, b starlark.StringDict) starlark.StringDict {
	for k, _ := range b {
		if _, ok := a[k]; ok {
			delete(a, k)
		}
	}
	return a
}

//recursively go through the output of starlark to convert them back into go
func recursiveFromValue(input interface{}) interface{} {
	switch input.(type) {
	case map[interface{}]interface{}:
		out := make(map[interface{}]interface{})
		for k, v := range input.(map[interface{}]interface{}) {
			out[k] = recursiveFromValue(v)
		}
		return out
	case []interface{}:
		out := input.([]interface{})
		for i := 0; i < len(out); i++ {
			out[i] = recursiveFromValue(out[i])
		}
		return out
	case starlark.Value:
		return convert.FromValue(input.(starlark.Value))
	default:
		return input
	}
}
