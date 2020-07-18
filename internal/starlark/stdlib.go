package starlark

import (
	"fmt"
	"github.com/chasehensel/starlight/convert"
	"go.starlark.net/starlark"
	"net/url"
	"reflect"
	"regexp"
)

var (
	ErrInvalidInput = fmt.Errorf("Bad input:")
)

//Handle the many repetitive errors gracefully
type errWriter struct {
	err error
}

func (ew *errWriter) assertType(val interface{}, t interface{}) interface{} {
	if ew.err != nil {
		return nil
	}
	if reflect.TypeOf(val) != reflect.TypeOf(t) {
		ew.err = fmt.Errorf("%w expected type %T, but found %T", ErrInvalidInput, t, val)
		return nil
	}
	return val
}

func (ew *errWriter) assertString(val interface{}) string {
	x := ew.assertType(val, "")
	if ew.err != nil {
		return ""
	}
	return x.(string)
}

func (ew *errWriter) assertInt64(val interface{}) int64 {
	var i int64 = 0
	x := ew.assertType(val, i)
	if ew.err != nil {
		return i
	}
	return x.(int64)
}

func compile(pattern interface{}) (*regexp.Regexp, error) {
	ew := errWriter{}
	patternS := ew.assertString(pattern)
	if ew.err != nil {
		return nil, ew.err
	}
	return regexp.Compile(patternS)
}

func match(regExp, match interface{}) (bool, error) {
	ew := errWriter{}
	var re *regexp.Regexp
	r := ew.assertType(regExp, re)
	matchS := ew.assertString(match)
	if ew.err != nil {
		return false, ew.err
	}
	return (r.(*regexp.Regexp)).MatchString(matchS), nil
}

type re struct {
	Match   func(regExp, match interface{}) (bool, error)
	Compile func(s interface{}) (*regexp.Regexp, error)
}

func StdLib(input starlark.Value, c *call) map[string]interface{} {
	env := map[string]interface{}{
		"args":  input,
		"error": func(str string, a ...interface{}) { c.err = fmt.Sprintf(str, a...) },
		"result": func(arg interface{}) {
			c.result = arg
		},
		//todo wrap this if/when it becomes useful
		"re": &re{Compile: compile, Match: match},
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
			ew := errWriter{}
			us := ew.assertString(u)
			if ew.err != nil {
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
	Env    map[string]interface{}
	result interface{}
	err    interface{}
	msgs   string
}

func CreateEnv(args starlark.Value, c *call) (starlark.StringDict, error) {
	stdlib := StdLib(args, c)
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
