package starlark

import (
	"awans.org/aft/internal/db"
	"fmt"
	"github.com/google/uuid"
	"github.com/starlight-go/starlight/convert"
	"go.starlark.net/starlark"
	"math"
	"net/http"
	"net/url"
	"regexp"
)

func ReplLib(tx db.RWTx) map[string]interface{} {
	env := map[string]interface{}{
		"findOne": func(modelName string, matcher db.Matcher) (r db.Record, err error) {
			model, err := tx.GetModel(modelName)
			if err != nil {
				return nil, err
			}
			return tx.FindOne(model.ID, matcher)
		},
		"findMany": func(modelName string, matcher db.Matcher) []db.Record {
			model, err := tx.GetModel(modelName)
			if err != nil {
				var results []db.Record
				return results
			}
			return tx.FindMany(model.ID, matcher)
		},
		"makeRecord": func(modelName string) db.Record {
			model, err := tx.GetModel(modelName)
			if err != nil {
				return nil
			}
			return tx.MakeRecord(model.ID)
		},
		"insert": func(r db.Record) {
			tx.Insert(r)
		},
		//Should we expose loadRel too?
		"connect": func(from, to db.Record, fromRel db.Relationship) {
			tx.Connect(from, to, fromRel)
		},
		"Eq": func(key, val string) db.Matcher {
			return db.Eq(key, val)
		},
		"EqFK": func(field, u string) db.Matcher {
			id, _ := uuid.Parse(u)
			return db.EqFK(field, id)
		},
		"And": func(matchers ...db.Matcher) db.Matcher {
			return db.And(matchers...)
		},
		"http": http.DefaultClient,
	}
	return env
}

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

func (s *StarlarkFunctionHandle) StdLib(input interface{}) map[string]interface{} {
	env := map[string]interface{}{

		"arg":    input,
		"error":  func(str string, a ...interface{}) { s.err = fmt.Sprintf(str, a...) },
		"print":  func(str string, a ...interface{}) { s.result = fmt.Sprintf(str, a...) },
		"re":     &re{Compile: compile, Match: match},
		"result": func(arg interface{}) { s.result = arg },
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

func (s *StarlarkFunctionHandle) createGlobals(args interface{}) (starlark.StringDict, error) {
	local := s.StdLib(args)
	globals, err := convert.MakeStringDict(local)
	if err != nil {
		return nil, err
	}
	api, err := convert.MakeStringDict(s.Env)
	if err != nil {
		return nil, err
	}
	//API overwrites local
	globals = clobber(globals, api)
	globals = union(globals, api)
	return globals, nil
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
