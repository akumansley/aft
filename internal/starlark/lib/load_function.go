package lib

import (
	"context"
	"errors"

	"awans.org/aft/internal/db"
	"awans.org/aft/internal/starlark/handlers"
	"go.starlark.net/starlark"
)

type ContextValue struct {
	context.Context
}

func (c ContextValue) String() string {
	return "context()"
}

func (c ContextValue) Type() string {
	return "context.Context"
}

func (c ContextValue) Freeze() {}

func (c ContextValue) Truth() starlark.Bool {
	return starlark.Bool(c.Context != nil)
}

func (c ContextValue) Hash() (uint32, error) {
	return 0, errors.New("Unhashable")
}

var loadFunction = starlark.NewBuiltin("loadFunction", loadFunctionFunc)

func loadFunctionFunc(thread *starlark.Thread, b *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	var ctx ContextValue
	var fname string

	if err := starlark.UnpackPositionalArgs(b.Name(), args, kwargs, 2, &ctx, &fname); err != nil {
		return starlark.None, err
	}

	rwtx, ok := db.RWTxFromContext(ctx)
	if !ok {
		return starlark.None, errors.New("loadFunction called in non-rwtx context")
	}
	functions := rwtx.Ref(db.FunctionInterface.ID())
	frec, err := rwtx.Query(functions, db.Filter(functions, db.Eq("name", fname))).OneRecord()
	if err != nil {
		return starlark.None, err
	}
	f, err := rwtx.Schema().LoadFunction(frec)
	if err != nil {
		return starlark.None, err
	}

	wrappedF := func(thread *starlark.Thread, fn *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
		// indirect
		var argSlice []interface{}
		for i := 0; i < f.Arity(); i++ {
			v := new(interface{})
			argSlice = append(argSlice, v)
		}

		if err := starlark.UnpackPositionalArgs(fname, args, kwargs, f.Arity(), argSlice...); err != nil {
			return starlark.None, err
		}
		// de-indirect
		var goArgs []interface{}
		for _, ptr := range argSlice {
			iv := *ptr.(*interface{})
			sv := iv.(starlark.Value)

			gv, err := handlers.Encode(sv)
			if err != nil {
				return starlark.None, err
			}
			goArgs = append(goArgs, gv)
		}

		rval, err := f.Call(goArgs)
		if err != nil {
			return starlark.None, err
		}
		decoded, err := handlers.Decode(rval)
		return decoded, err
	}

	fv := starlark.NewBuiltin(fname, wrappedF)
	return fv, nil
}
