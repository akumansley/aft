package handlers

import (
	"awans.org/aft/internal/db"
	"fmt"
	"github.com/chasehensel/starlight/convert"
	"go.starlark.net/starlark"
	"go.starlark.net/starlarkstruct"
)

type Handler struct {
	tx db.RWTx
}

func API(tx db.RWTx) *starlarkstruct.Module {
	h := &Handler{tx: tx}
	return &starlarkstruct.Module{
		Name: "aft",
		Members: starlark.StringDict{
			"findOne":    starlark.NewBuiltin("aft.findOne", h.findOne),
			"findMany":   starlark.NewBuiltin("aft.findMany", h.findMany),
			"delete":     starlark.NewBuiltin("aft.delete", h.del),
			"deleteMany": starlark.NewBuiltin("aft.deleteMany", h.deleteMany),
			"update":     starlark.NewBuiltin("aft.update", h.update),
			"updateMany": starlark.NewBuiltin("aft.updateMany", h.updateMany),
			"upsert":     starlark.NewBuiltin("aft.upsert", h.upsert),
			"create":     starlark.NewBuiltin("aft.create", h.create),
			"count":      starlark.NewBuiltin("aft.count", h.count),
		},
	}
}

func unpack(thread *starlark.Thread, b *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (modelName string, body map[string]interface{}, err error) {
	var starlarkModelName starlark.Value
	var starlarkBody starlark.Value
	if err := starlark.UnpackPositionalArgs(b.Name(), args, kwargs, 2, &starlarkModelName, &starlarkBody); err != nil {
		return modelName, body, err
	}

	modelName, ok := starlark.AsString(starlarkModelName)
	if !ok {
		return modelName, body, fmt.Errorf("Invalid model: %s", starlarkModelName)
	}
	inp := convert.FromValue(starlarkBody)
	ibody, err := encode(inp)
	if err != nil {
		return modelName, body, err
	}
	body, ok = ibody.(map[string]interface{})
	if !ok {
		return modelName, body, fmt.Errorf("Invalid arguments: %s", ibody)
	}
	return modelName, body, nil
}

//recursively go through starlark to convert them back into go types
func encode(input interface{}) (interface{}, error) {
	switch input.(type) {
	case map[interface{}]interface{}:
		out := make(map[string]interface{})
		for k, v := range input.(map[interface{}]interface{}) {
			if key, ok := k.(string); ok {
				enc, err := encode(v)
				if err != nil {
					return nil, err
				}
				out[key] = enc
			} else {
				return nil, fmt.Errorf("Key %+v is type %T, not string.", k, k)
			}
		}
		return out, nil
	case []interface{}:
		out := input.([]interface{})
		for i := 0; i < len(out); i++ {
			enc, err := encode(out[i])
			if err != nil {
				return nil, err
			}
			out[i] = enc
		}
		return out, nil
	case starlark.Value:
		return convert.FromValue(input.(starlark.Value)), nil
	default:
		return input, nil
	}
}

func output(result interface{}) (starlark.Value, error) {
	if count, ok := result.(int); ok {
		return starlark.MakeInt(count), nil
	}
	return decode(result)
}

//recursively go through starlark to convert them back into go
func decode(input interface{}) (starlark.Value, error) {
	switch input.(type) {
	case *db.QueryResult:
		rec, _ := input.(*db.QueryResult)
		if rec == nil {
			return starlark.None, nil
		}
		return decodeQR(rec)
	case []*db.QueryResult:
		recs, _ := input.([]*db.QueryResult)
		var outs starlark.Tuple
		for _, rec := range recs {
			val, err := decodeQR(rec)
			if err != nil {
				return starlark.None, err
			}
			outs = append(outs, val)
		}
		return outs, nil
	}
	panic("Unknown type in decoding query results")
}

func decodeQR(rec *db.QueryResult) (starlark.Value, error) {
	m, err := rec.Map()
	if err != nil {
		return starlark.None, err
	}
	for k, v := range rec.ToOne {
		out, err := decode(v)
		if err != nil {
			return starlark.None, err
		}
		m[k] = out
	}
	for k, v := range rec.ToMany {
		out, err := decode(v)
		if err != nil {
			return starlark.None, err
		}
		m[k] = out
	}

	sd := make(starlark.StringDict)
	for k, v := range m {
		val, err := convert.ToValue(v)
		if err != nil {
			return starlark.None, err
		}
		sd[k] = val
	}
	return &starlarkstruct.Module{
		Name:    "record",
		Members: sd,
	}, err
}
