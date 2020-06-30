package starlark

import (
	"awans.org/aft/internal/datatypes"
	"awans.org/aft/internal/db"
	"fmt"
	"github.com/google/uuid"
	"go.starlark.net/starlark"
	"reflect"
)

var (
	ErrInvalidInput = fmt.Errorf("Bad input:")
)

//Handle the many repetitive errors gracefully
type errWriter struct {
	err error
}

func (ew *errWriter) assertType(val interface{}, t interface{}) interface{} {
	if ew.err != nil {
		return nil
	}
	if reflect.TypeOf(val) != reflect.TypeOf(t) {
		ew.err = fmt.Errorf("%w expected type %T, but found %T", ErrInvalidInput, t, val)
		return nil
	}
	return val
}

func (ew *errWriter) assertString(val interface{}) string {
	x := ew.assertType(val, "")
	if ew.err != nil {
		return ""
	}
	return x.(string)
}

func (ew *errWriter) assertUUID(val interface{}) uuid.UUID {
	u := uuid.UUID{}
	x := ew.assertType(val, u)
	if ew.err != nil {
		return uuid.Nil
	}
	return x.(uuid.UUID)
}

func (ew *errWriter) assertID(val interface{}) db.ID {
	u := db.ID(uuid.UUID{})
	x := ew.assertType(val, u)
	if ew.err != nil {
		return u
	}
	return x.(db.ID)
}

func (ew *errWriter) assertInt64(val interface{}) int64 {
	var i int64 = 0
	x := ew.assertType(val, i)
	if ew.err != nil {
		return i
	}
	return x.(int64)
}

func (ew *errWriter) assertModel(val interface{}, tx db.RWTx) db.Model {
	name := ew.assertString(val)
	if ew.err != nil {
		return db.Model{}
	}
	m, err := tx.GetModel(name)
	if err != nil {
		ew.err = err
		return db.Model{}
	}
	return m
}

func (ew *errWriter) assertMatcher(val interface{}) db.Matcher {
	if val, ok := val.(db.Matcher); ok {
		return val
	}
	ew.err = fmt.Errorf("%w %T doesn't implement Matcher interface", ErrInvalidInput, val)
	return db.FieldMatcher{}
}

func (ew *errWriter) assertMap(val interface{}) map[interface{}]interface{} {
	empty := make(map[interface{}]interface{})
	ma := ew.assertType(val, empty)
	if ew.err != nil {
		return empty
	}
	return ma.(map[interface{}]interface{})
}

func (ew *errWriter) assertStarlarkRecord(val interface{}) *starlarkRecord {
	r := &starlarkRecord{}
	out := ew.assertType(val, r)
	if ew.err != nil {
		return r
	}
	return out.(*starlarkRecord)
}

func (ew *errWriter) GetFromRecord(s string, r Record) interface{} {
	if ew.err != nil {
		return nil
	}
	out, err := r.Get(s)
	if err != nil {
		ew.err = err
		return nil
	}
	return out
}

func (ew *errWriter) SetDBRecord(s string, i interface{}, r db.Record) {
	if ew.err != nil {
		return
	}
	err := r.Set(s, i)
	if err != nil {
		ew.err = err
	}
}

//Wrapper for the Record interface so we can control which methods to expose.
// This gets surfaced in Starlark as return values of database functions
type Record interface {
	Map() map[string]interface{}
	ID() db.ID
	Get(string) (interface{}, error)
	GetFK(string) (db.ID, error)
}

type starlarkRecord struct {
	inner db.Record
}

func (r *starlarkRecord) Map() map[string]interface{} {
	return r.inner.Map()
}

func (r *starlarkRecord) ID() db.ID {
	return r.inner.ID()
}

func (r *starlarkRecord) Get(fieldName string) (interface{}, error) {
	field, err := r.inner.Get(fieldName)
	if err != nil {
		return nil, err
	}
	return field, nil
}

func (r *starlarkRecord) GetFK(fieldName string) (db.ID, error) {
	rel, err := r.inner.GetFK(fieldName)
	if err != nil {
		return db.ID(uuid.Nil), err
	}
	return rel, nil
}

//Actual DB API
func DBLib(tx db.RWTx) map[string]interface{} {
	env := make(map[string]interface{})
	env["FindOne"] = func(mn, mm interface{}) (Record, error) {
		ew := errWriter{}
		m := ew.assertModel(mn, tx)
		ma := ew.assertMatcher(mm)
		if ew.err != nil {
			return nil, ew.err
		}
		r, err := tx.FindOne(m.ID, ma)
		if err != nil {
			return nil, err
		}
		return &starlarkRecord{inner: r}, nil
	}
	env["FindMany"] = func(mn, mm interface{}) ([]Record, error) {
		ew := errWriter{}
		m := ew.assertModel(mn, tx)
		ma := ew.assertMatcher(mm)
		if ew.err != nil {
			return nil, ew.err
		}
		recs, err := tx.FindMany(m.ID, ma)
		if err != nil {
			return nil, err
		}
		var out []Record
		for i := 0; i < len(recs); i++ {
			out = append(out, &starlarkRecord{inner: recs[i]})
		}
		return out, nil
	}
	env["Eq"] = func(k, v interface{}) (db.Matcher, error) {
		ew := errWriter{}
		key := ew.assertString(k)
		if ew.err != nil {
			return nil, ew.err
		}
		return db.Eq(key, v), nil
	}
	env["EqID"] = func(v interface{}) (db.Matcher, error) {
		ew := errWriter{}
		id := ew.assertID(v)
		if ew.err != nil {
			return nil, ew.err
		}
		return db.EqID(id), nil
	}
	env["ID"] = func(v interface{}) (db.ID, error) {
		ew := errWriter{}
		id := ew.assertUUID(v)
		if ew.err != nil {
			return db.ID(uuid.Nil), ew.err
		}
		return db.ID(id), nil
	}
	env["EqFK"] = func(k, v interface{}) (db.Matcher, error) {
		ew := errWriter{}
		key := ew.assertString(k)
		id := ew.assertID(v)
		if ew.err != nil {
			return nil, ew.err
		}
		return db.EqFK(key, id), nil
	}
	env["And"] = func(matchers ...interface{}) (db.Matcher, error) {
		ew := errWriter{}
		var out []db.Matcher
		for i := 0; i < len(matchers); i++ {
			m := ew.assertMatcher(matchers[i])
			if ew.err != nil {
				return nil, ew.err
			}
			out = append(out, m)
		}
		return db.And(out...), nil
	}
	env["Insert"] = func(mn interface{}, fields interface{}) (Record, error) {
		ew := errWriter{}
		m := ew.assertModel(mn, tx)
		r, err := tx.MakeRecord(m.ID)
		if err != nil {
			return nil, err
		}
		ew.SetDBRecord("id", uuid.New(), r)
		fieldMap := ew.assertMap(fields)
		for key, val := range fieldMap {
			ks := ew.assertString(key)
			ew.SetDBRecord(ks, recursiveFromValue(val.(starlark.Value)), r)
		}
		if ew.err != nil {
			return nil, ew.err
		}
		tx.Insert(r)
		return &starlarkRecord{inner: r}, nil
	}
	env["Update"] = func(r interface{}, fields interface{}) (Record, error) {
		ew := errWriter{}
		rec := ew.assertStarlarkRecord(r)
		if ew.err != nil {
			return nil, ew.err
		}
		oldRec := rec.inner
		newRec := oldRec.DeepCopy()
		fieldMap := ew.assertMap(fields)
		for key, val := range fieldMap {
			ks := ew.assertString(key)
			ew.SetDBRecord(ks, recursiveFromValue(val.(starlark.Value)), newRec)
		}
		if ew.err != nil {
			return nil, ew.err
		}
		err := tx.Update(oldRec, newRec)
		if err != nil {
			return nil, err
		}
		return &starlarkRecord{inner: newRec}, err

	}
	env["Connect"] = func(s interface{}, r1 interface{}, r2 interface{}) (bool, error) {
		ew := errWriter{}
		bname := ew.assertString(s)
		rec1 := ew.assertStarlarkRecord(r1)
		rec2 := ew.assertStarlarkRecord(r2)
		if ew.err != nil {
			return false, ew.err
		}
		binding, err := rec1.inner.Model().GetBinding(bname)
		if err != nil {
			return false, err
		}
		err = tx.Connect(rec1.inner, rec2.inner, binding.Relationship)
		if err != nil {
			return false, err
		}
		return true, nil
	}
	env["Delete"] = func(r interface{}) (Record, error) {
		ew := errWriter{}
		rec := ew.assertStarlarkRecord(r)
		if ew.err != nil {
			return nil, ew.err
		}
		err := tx.Delete(rec.inner)
		if err != nil {
			return nil, err
		}
		return rec, err
	}
	env["Exec"] = func(r interface{}, args interface{}) (interface{}, error) {
		ew := errWriter{}
		rec := ew.assertStarlarkRecord(r)
		ci := ew.GetFromRecord("code", rec)
		code := ew.assertString(ci)
		runtime, err := db.RecordToEnumValue(rec.inner, "runtime", tx)
		if err != nil {
			return nil, err
		}
		fs, err := db.RecordToEnumValue(rec.inner, "functionSignature", tx)
		if err != nil {
			return nil, ew.err
		}
		switch runtime.ID {
		case db.Starlark.ID:
			fh := &StarlarkFunctionHandle{Code: code, FunctionSignature: db.FunctionSignatureEnumValue{fs}}
			return fh.Invoke(args)
		case db.Native.ID:
			fh := &datatypes.GoFunctionHandle{Function: db.CodeMap[rec.ID()].Function}
			return fh.Invoke(args)
		default:
			return nil, fmt.Errorf("Can't execute code because it uses an unsupported runtime")
		}
	}
	return env
}
