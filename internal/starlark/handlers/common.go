package handlers

import (
	"encoding/json"
	"errors"
	"fmt"

	"awans.org/aft/internal/db"
	"github.com/chasehensel/starlight/convert"
	"go.starlark.net/starlark"
)

var NonStringKey = errors.New("starlark.Dict had a non-string key")

type qr struct {
	*db.QueryResult
}

func (q *qr) MarhsalJSON() ([]byte, error) {
	return q.QueryResult.MarshalJSON()
}

func (q *qr) Attr(name string) (starlark.Value, error) {
	val, err := q.QueryResult.Get(name)
	if err != nil {
		return starlark.None, err
	}
	return Decode(val)
}

func (q *qr) AttrNames() (attrNames []string) {
	aMap, err := q.QueryResult.Map()
	if err != nil {
		panic(err)
	}
	keys := make([]string, len(aMap))

	i := 0
	for k := range aMap {
		keys[i] = k
		i++
	}

	return keys
}

func (q *qr) String() string {
	return q.QueryResult.String()
}

func (q *qr) Type() string {
	return "queryresult"
}

func (q *qr) Freeze() {

}

func (q *qr) Truth() starlark.Bool {
	return true
}

func (q *qr) Hash() (uint32, error) {
	return 0, errors.New("queryresult is unhashable")
}

type srec struct {
	db.Record
}

func (s *srec) MarshalJSON() ([]byte, error) {
	bytes, err := json.Marshal(s.Record)
	return bytes, err
}

func (s *srec) Attr(name string) (starlark.Value, error) {
	val, err := s.Record.Get(name)
	if err != nil {
		return starlark.None, err
	}
	return Decode(val)
}

func (s *srec) AttrNames() (attrNames []string) {
	return s.Record.FieldNames()
}

func (s *srec) String() string {
	return fmt.Sprintf("%v", s.Record.Map())
}

func (s *srec) Type() string {
	return s.Record.Type()
}

func (s *srec) Freeze() {

}

func (s *srec) Truth() starlark.Bool {
	return true
}

func (s *srec) Hash() (uint32, error) {
	return 0, errors.New("record is unhashable")
}

func tryJsonDict(m *starlark.Dict) (map[string]interface{}, error) {
	out := make(map[string]interface{})
	for _, k := range m.Keys() {
		key, ok := starlark.AsString(k)
		if !ok {
			return nil, NonStringKey
		}
		val, _, _ := m.Get(k)
		var err error
		out[key], err = Encode(val)
		if err != nil {
			return nil, err
		}
	}
	return out, nil
}

//recursively go through starlark to convert them back into go types
func Encode(input interface{}) (interface{}, error) {
	switch input.(type) {
	case *starlark.Dict:
		out, err := tryJsonDict(input.(*starlark.Dict))
		if err == NonStringKey {
			return convert.FromDict(input.(*starlark.Dict)), nil
		}
		if err != nil {
			return nil, err
		}
		return out, nil
	case map[interface{}]interface{}:
		out := make(map[string]interface{})
		for k, v := range input.(map[interface{}]interface{}) {
			if key, ok := k.(string); ok {
				enc, err := Encode(v)
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
			enc, err := Encode(out[i])
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
	return Decode(result)
}

//recursively go through go values to convert them into starlark
func Decode(input interface{}) (starlark.Value, error) {
	if input == nil {
		return starlark.None, nil
	}

	switch input.(type) {
	case db.Record:
		r := input.(db.Record)
		return &srec{r}, nil
	case *db.QueryResult:
		rec := input.(*db.QueryResult)
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
	default:
		return convert.ToValue(input)
	}
}

func decodeQR(rec *db.QueryResult) (starlark.Value, error) {
	return &qr{rec}, nil
}
