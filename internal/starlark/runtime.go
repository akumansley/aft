package starlark

import (
	"context"
	"fmt"
	"io/ioutil"

	"awans.org/aft/internal/db"
	"awans.org/aft/internal/errors"
	"awans.org/aft/internal/starlark/lib"
	"github.com/chasehensel/starlight/convert"
	"github.com/markbates/pkger"
	"go.starlark.net/resolve"
	"go.starlark.net/starlark"
)

func loadCode(path string) []byte {
	f, err := pkger.Open(path)
	if err != nil {
		panic(err)
	}
	b, err := ioutil.ReadAll(f)
	if err != nil {
		panic(err)
	}
	return b
}

// Loader

func NewStarlarkRuntime() *StarlarkRuntime {
	return &StarlarkRuntime{}
}

type StarlarkRuntime struct {
}

//configure starlark
func init() {
	resolve.AllowNestedDef = true // allow nested def
	resolve.AllowFloat = true     // allow floating point literals, the 'float' built-in, and x / y
	resolve.AllowSet = true       // allow the 'set' built-in
	resolve.AllowRecursion = true // allow while statements and recursive functions
}

type call struct {
	msgs []string
}

func prepareInput(thread *starlark.Thread, args []interface{}) (vals []starlark.Value, err error) {
	for _, arg := range args {
		var val starlark.Value
		val, err = lib.ToStarlark(arg)
		if err != nil {
			return
		}
		vals = append(vals, val)
	}
	return
}

func (sr *StarlarkRuntime) Execute(ctx context.Context, code string, args []interface{}) (interface{}, error) {
	c := &call{}

	thread := &starlark.Thread{
		Load: nil,
		Print: func(_ *starlark.Thread, msg string) {
			c.msgs = append(c.msgs, msg)
			fmt.Println(msg)
		},
	}
	thread.SetLocal("ctx", ctx)

	convertedArgs, err := prepareInput(thread, args)
	if err != nil {
		return nil, err
	}

	env := lib.Lib(ctx)
	globals, err := starlark.ExecFile(thread, "", []byte(code), env)
	if err != nil {
		return nil, err
	}

	if globals["main"] == nil {
		return nil, fmt.Errorf("Missing main function")
	}
	out, err := starlark.Call(thread, globals["main"], convertedArgs, nil)
	if err != nil {
		return nil, err
	}
	if serr, ok := out.(*lib.StarlarkError); ok {
		return nil, errors.AftError{Code: serr.Code, Message: serr.Message}
	}
	goland := recursiveFromValue(out)
	return goland, nil
}

func (sr *StarlarkRuntime) ProvideModel() db.ModelL {
	return StarlarkFunctionModel
}

func (sr *StarlarkRuntime) Load(rec db.Record) db.Function {
	return &starlarkFunction{rec, sr}
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
