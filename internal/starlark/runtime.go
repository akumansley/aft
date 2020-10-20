package starlark

import (
	"context"
	"fmt"
	"io/ioutil"

	"awans.org/aft/internal/db"
	"awans.org/aft/internal/starlark/lib"
	"github.com/chasehensel/starlight/convert"
	"github.com/markbates/pkger"
	"go.starlark.net/resolve"
	"go.starlark.net/starlark"
)

func init() {
	pkger.Include("/internal/starlark/aft.star")
}

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

func prepareAft(ctx context.Context, thread *starlark.Thread) (aft starlark.Value, err error) {
	aftCode := loadCode("/internal/starlark/aft.star")
	libVars, err := starlark.ExecFile(thread, "", aftCode, lib.Lib)

	if err != nil {
		return nil, err
	}

	args := []starlark.Value{
		lib.ContextValue{Context: ctx},
	}

	return starlark.Call(thread, libVars["preamble"], args, nil)
}

func prepareInput(thread *starlark.Thread, args []interface{}) (vals []starlark.Value, err error) {
	for _, arg := range args {
		var val starlark.Value
		if ctx, ok := arg.(context.Context); ok {
			val, err = prepareAft(ctx, thread)
			if err != nil {
				return
			}
		} else {
			val, err = convert.ToValue(arg)
			if err != nil {
				return
			}
		}
		vals = append(vals, val)
	}
	return
}

func (sr *StarlarkRuntime) Execute(code string, args []interface{}) (interface{}, error) {
	c := &call{}

	thread := &starlark.Thread{
		Load: nil,
		Print: func(_ *starlark.Thread, msg string) {
			c.msgs = append(c.msgs, msg)
			fmt.Println(msg)
		},
	}

	convertedArgs, err := prepareInput(thread, args)
	if err != nil {
		return nil, err
	}

	globals, err := starlark.ExecFile(thread, "", []byte(code), lib.Lib)
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
	goland := recursiveFromValue(out)
	return goland, nil
}

func (sr *StarlarkRuntime) ProvideModel() db.ModelL {
	return StarlarkFunctionModel
}

func (sr *StarlarkRuntime) Load(tx db.Tx, rec db.Record) db.Function {
	return &starlarkFunction{rec, sr, tx}
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
