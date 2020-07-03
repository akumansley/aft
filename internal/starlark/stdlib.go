package starlark

import (
	"fmt"
	"github.com/chasehensel/starlight/convert"
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

func (s *StarlarkFunctionHandle) StdLib(input starlark.Value) map[string]interface{} {
	env := map[string]interface{}{

		"args":   input,
		"test":   func(str string, a ...interface{}) { fmt.Printf(str, a...) },
		"error":  func(str string, a ...interface{}) { s.err = fmt.Sprintf(str, a...) },
		"print":  func(str string, a ...interface{}) { s.result = fmt.Sprintf(str, a...) },
		"re":     &re{Compile: compile, Match: match},
		"result": func(arg interface{}) { s.result = arg },
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
		if s.Env != nil {
			a, amax := PrintPretty(env, 0)
			b, bmax := PrintPretty(s.Env, 0)
			if amax > bmax {
				b, bmax = PrintPretty(s.Env, amax)
			} else {
				a, amax = PrintPretty(env, bmax)
			}
			s.result = a + b
		} else {
			a, _ := PrintPretty(env, 0)
			s.result = a
		}

	}
	return env
}

func (s *StarlarkFunctionHandle) createEnv(args starlark.Value) (starlark.StringDict, error) {
	stdlib := s.StdLib(args)
	env, err := convert.MakeStringDict(stdlib)
	if err != nil {
		return nil, err
	}
	api, err := convert.MakeStringDict(s.Env)
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
