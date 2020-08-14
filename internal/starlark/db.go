package starlark

import (
	"awans.org/aft/internal/db"
	"awans.org/aft/internal/starlark/handlers"
	"fmt"
	"github.com/chasehensel/starlight/convert"
	"go.starlark.net/resolve"
	"go.starlark.net/starlark"
	"go.starlark.net/starlarkstruct"
	"go.starlark.net/syntax"
)

func createFunctionDict(tx db.RWTx) (starlark.StringDict) {
	sd := make(starlark.StringDict)
	for _, c := range tx.Query(tx.Ref(db.FunctionInterface.ID())).All() {
		code := c.Record
		name := code.MustGet("name").(string)
		fn, _ := tx.Schema().LoadFunction(code)
		var wrapper = func(i interface{}) (starlark.Value, starlark.Value) {
			out, err := fn.Call(i)
			if err != nil {
				return starlark.None, starlark.String(fmt.Sprintf("%s", err))
			}
			outS, _ := convert.ToValue(out)
			return outS, starlark.None
		}
		sf, _ := convert.ToValue(wrapper)
		sd[name] = sf
	}
	return sd
}

func getFunction(sd starlark.StringDict) (func(thread *starlark.Thread, b *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (out starlark.Value, err error)) {
	return func(thread *starlark.Thread, b *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (out starlark.Value, err error) {
		var val starlark.Value
		if err := starlark.UnpackPositionalArgs(b.Name(), args, kwargs, 1, &val); err != nil {
			return starlark.None, err
		}
		s, ok := starlark.AsString(val)
		if !ok {
			return starlark.None, fmt.Errorf("Invalid string: %s", val)
		}
		if _, ok := sd[s]; !ok {
			return starlark.None, nil
		}
		return sd[s], nil
	}
}

func DBLib(tx db.RWTx) map[string]interface{} {
	env := make(map[string]interface{})
	sd := createFunctionDict(tx)
	env["aft"] = &starlarkstruct.Module{
		Name: "aft",
		Members: starlark.StringDict{
			"api": handlers.API(tx),
			"function": &starlarkstruct.Module{
				Name:    "function",
				Members: sd,
			},
			"getFunction": starlark.NewBuiltin("getFunction", getFunction(sd)),
		},
	}
	env["parse"] = func(code interface{}) (string, bool, error) {
		if input, ok := code.(string); ok {
			f, err := syntax.Parse("", input, 0)
			if err != nil {
				return fmt.Sprintf("%s", err), false, nil
			}
			var isPredeclared = func(s string) bool {
				c := &call{}
				c.Env = DBLib(tx)
				env, err := CreateEnv(c)
				if err != nil {
					return false
				}

				if _, ok := env[s]; ok {
					return true
				}
				if _, ok := StdLib(nil)[s]; ok {
					return true
				}
				return false
			}
			err = resolve.File(f, isPredeclared, starlark.Universe.Has)
			if err != nil {
				return fmt.Sprintf("%s", err), false, nil
			}
			return "", true, nil
		}
		return "", false, fmt.Errorf("%w code was type %T", ErrInvalidInput, code)
	}
	env["exec"] = func(code interface{}, args interface{}) (string, bool, error) {
		if input, ok := code.(string); ok {
			sh := MakeStarlarkFunction(db.NewID(), "", db.RPC, input)
			env := DBLib(tx)
			r, err := sh.CallWithEnv(args, env)
			if err != nil {
				if evale, ok := err.(*starlark.EvalError); ok {
					return evale.Backtrace(), false, nil
				}
				return fmt.Sprintf("%s", err), false, nil
			}
			if r == nil {
				return "", true, nil
			}
			return fmt.Sprintf("%v", r), true, nil
		}
		return "", false, fmt.Errorf("%w code was type %T", ErrInvalidInput, code)
	}
	return env
}
