package lib

import (
	"context"
	"errors"

	"awans.org/aft/internal/db"
	"go.starlark.net/starlark"
	"go.starlark.net/starlarkstruct"
)

func Lib(ctx context.Context) starlark.StringDict {
	lib := starlark.StringDict{}

	for k, v := range staticLib {
		lib[k] = v
	}
	lib["func"] = &magicFunc{ctx}
	return lib
}

var staticLib = starlark.StringDict{
	"re":         re,
	"asSelf":     asSelf,
	"asUser":     asUser,
	"urlparse":   urlparse,
	"struct":     starlark.NewBuiltin("struct", starlarkstruct.Make),
	"error":      starlark.NewBuiltin("error", makeError),
	"findOne":    makeBuiltin("findOne", 2),
	"findMany":   makeBuiltin("findMany", 2),
	"count":      makeBuiltin("count", 2),
	"upsert":     makeBuiltin("upsert", 2),
	"update":     makeBuiltin("update", 2),
	"updateMany": makeBuiltin("updateMany", 2),
	"create":     makeBuiltin("create", 2),
	"delete":     makeBuiltin("delete", 2),
	"deleteMany": makeBuiltin("deleteMany", 2),
}

func makeBuiltin(name string, arity int) starlark.Value {
	wrapperFunc := func(thread *starlark.Thread, fn *starlark.Builtin,
		args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
		ctx := thread.Local("ctx").(context.Context)
		tx, ok := db.TxFromContext(ctx)
		if !ok {
			return starlark.None, errors.New("No tx in starlark call")
		}

		// indirect
		var argSlice []interface{}
		for i := 0; i < arity; i++ {
			v := new(interface{})
			argSlice = append(argSlice, v)
		}

		if err := starlark.UnpackPositionalArgs(name, args, kwargs, arity, argSlice...); err != nil {
			return starlark.None, err
		}

		// de-indirect
		var goArgs []interface{}
		for _, ptr := range argSlice {
			iv := *ptr.(*interface{})
			sv := iv.(starlark.Value)

			gv, err := FromStarlark(sv)
			if err != nil {
				return starlark.None, err
			}
			goArgs = append(goArgs, gv)
		}

		// TODO what problems does it cause to allow functions to call functions
		f, err := tx.Schema().GetFunctionByName(name)
		if err != nil {
			return starlark.None, err
		}
		rval, err := f.Call(tx.Context(), goArgs)
		if err != nil {
			return starlark.None, err
		}

		decoded, err := ToStarlark(rval)
		return decoded, err
	}
	return starlark.NewBuiltin(name, wrapperFunc)
}
