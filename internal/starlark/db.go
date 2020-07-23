package starlark

import (
	"awans.org/aft/internal/db"
	"fmt"
	"github.com/google/uuid"
	"go.starlark.net/resolve"
	"go.starlark.net/starlark"
	"go.starlark.net/syntax"
)

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

func (ew *errWriter) assertModel(val interface{}, tx db.RWTx) db.Model {
	name := ew.assertString(val)
	if ew.err != nil {
		return nil
	}
	m, err := tx.Schema().GetModel(name)
	if err != nil {
		ew.err = err
		return nil
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
	ID() db.ID
	Get(string) (interface{}, error)
}

type starlarkRecord struct {
	inner db.Record
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
		t := tx.Ref(m.ID())
		r, err := tx.Query(t).Filter(t, ma).OneRecord()
		return &starlarkRecord{inner: r}, err
	}
	env["FindMany"] = func(mn, mm interface{}) ([]Record, error) {
		ew := errWriter{}
		m := ew.assertModel(mn, tx)
		ma := ew.assertMatcher(mm)
		if ew.err != nil {
			return nil, ew.err
		}
		t := tx.Ref(m.ID())
		results := tx.Query(t).Filter(t, ma).All()
		var out []Record
		for _, res := range results {
			out = append(out, &starlarkRecord{inner: res.Record})
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
		r, err := tx.MakeRecord(m.ID())
		if err != nil {
			return nil, err
		}
		ew.SetDBRecord("id", uuid.New(), r)
		fieldMap := ew.assertMap(fields)
		for key, val := range fieldMap {
			ks := ew.assertString(key)
			ew.SetDBRecord(ks, val, r)
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
			ew.SetDBRecord(ks, val, newRec)
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

		var relationship db.Relationship
		rels, err := rec1.inner.Interface().Relationships()
		if err != nil {
			return false, err
		}
		found := false
		for _, rel := range rels {
			if rel.Name() == bname {
				relationship = rel
				found = true
				break
			}
		}
		if !found {
			return false, err
		}

		err = tx.Connect(rec1.inner.ID(), rec2.inner.ID(), relationship.ID())
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
	env["Parse"] = func(code interface{}) (string, bool, error) {
		if input, ok := code.(string); ok {
			f, err := syntax.Parse("", input, 0)
			if err != nil {
				return fmt.Sprintf("%s", err), false, nil
			}
			var isPredeclared = func(s string) bool {
				env, err := CreateEnv(starlark.String(""), nil)
				if err != nil {
					return false
				}

				if _, ok := env[s]; ok {
					return true
				}
				if _, ok := StdLib(starlark.String(""), nil)[s]; ok {
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
	env["Exec"] = func(code interface{}, args interface{}) (string, bool, error) {
		if rec, ok := code.(*starlarkRecord); ok {
			f, err := tx.Schema().LoadFunction(rec.inner)
			r, err := f.Call(args)
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
		} else if input, ok := code.(string); ok {
			sh := MakeStarlarkFunction(db.NewID(), "", db.RPC, input)
			r, err := sh.Call(args)
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
