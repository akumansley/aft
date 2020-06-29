package starlark

import (
	"fmt"
	"github.com/starlight-go/starlight/convert"
	"go.starlark.net/starlark"
	"math"
	"net/url"
	"regexp"
)

func compile(pattern string) *regexp.Regexp {
	return regexp.MustCompile(pattern)
}

func match(r *regexp.Regexp, match string) bool {
	return r.MatchString(match)
}

type re struct {
	Match   func(r *regexp.Regexp, match string) bool
	Compile func(pattern string) *regexp.Regexp
}

func StdLib(input interface{}, c *call) map[string]interface{} {
	env := map[string]interface{}{

		"args":  input,
		"test":  func(str string, a ...interface{}) { fmt.Printf(str, a...) },
		"error": func(str string, a ...interface{}) { c.err = fmt.Sprintf(str, a...) },
		"print": func(str string, a ...interface{}) { c.result = fmt.Sprintf(str, a...) },
		"re":    &re{Compile: compile, Match: match},
		"result": func(arg interface{}) {
			c.result = arg
		},
		//todo wrap this if/when it becomes useful
		"urlparse": func(us string) (*url.URL, bool) {
			u, e := url.ParseRequestURI(us)
			if e != nil {
				return nil, false
			}
			return u, true
		},
		"sprint": func(str string, a ...interface{}) string { return fmt.Sprintf(str, a...) },
	}
	env["man"] = func() {
		if c.Env != nil {
			a, amax := PrintPretty(env, 0)
			b, bmax := PrintPretty(c.Env, 0)
			if amax > bmax {
				b, bmax = PrintPretty(c.Env, amax)
			} else {
				a, amax = PrintPretty(env, bmax)
			}
			c.result = a + b
		} else {
			a, _ := PrintPretty(env, 0)
			c.result = a
		}

	}
	return env
}

type call struct {
	Env    map[string]interface{}
	result interface{}
	err    interface{}
}

func CreateEnv(args interface{}, c *call) (starlark.StringDict, error) {
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

func PrintPretty(input map[string]interface{}, existingMax int) (string, int) {
	max := existingMax - 2
	out := ""
	for k, _ := range input {
		//4 because python is tabbed four spaces in the repl
		cur := int(math.Floor(float64(len(k)) / 4.0))
		if cur > max {
			max = cur
		}
	}
	//add two offset
	max += 2
	for k, v := range input {
		//4 because python is tabbed four spaces in the repl
		cur := int(math.Floor(float64(len(k)) / 4.0))
		out = fmt.Sprintf("%s%s", out, k)

		for i := 0; i < max-cur; i++ {
			out = fmt.Sprintf("%s\t", out)
		}
		out = fmt.Sprintf("%s%T\n", out, v)
	}
	return out, max
}

//recursively go through the output of starlark to convert them back into go
func recursiveFromValue(input interface{}) interface{} {
	switch input.(type) {
	case *starlark.Dict:
		out := make(map[interface{}]interface{})
		m := input.(*starlark.Dict)
		for _, k := range m.Keys() {
			key := convert.FromValue(k)
			val, _, _ := m.Get(k)
			out[key] = recursiveFromValue(val)
		}
		return out
	case map[interface{}]interface{}:
		out := make(map[interface{}]interface{})
		for k, v := range input.(map[interface{}]interface{}) {
			out[k] = recursiveFromValue(v)
		}
		return out
	case *starlark.List:
		l := input.(*starlark.List)
		out := make([]interface{}, 0, l.Len())
		var v starlark.Value
		i := l.Iterate()
		defer i.Done()
		for i.Next(&v) {
			val := recursiveFromValue(v)
			out = append(out, val)
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

//recursively go through the input and convert into starlark
func recursiveToValue(input interface{}) (out interface{}, err error) {
	if err != nil {
		return nil, err
	}
	switch input.(type) {
	case map[interface{}]interface{}:
		out := make(map[interface{}]interface{})
		for k, v := range input.(map[interface{}]interface{}) {
			val, err := recursiveToValue(v)
			if err != nil {
				return nil, err
			}
			out[k] = val
		}
		return out, nil
	case []interface{}:
		out := input.([]interface{})
		for i := 0; i < len(out); i++ {
			val, err := recursiveToValue(out[i])
			if err != nil {
				return nil, err
			}
			out[i] = val
		}
		return out, nil
	default:
		return convert.ToValue(input)

	}
}
