package starlark

import (
	"awans.org/aft/internal/db"
	"fmt"
	"github.com/google/uuid"
	"go.starlark.net/starlark"
)

var (
	ErrInvalidInput = fmt.Errorf("Bad input:")
)

//Helper methods to make the API simpler
func getString(s interface{}) (string, error) {
	if val, ok := s.(string); ok {
		return val, nil
	}
	return "", fmt.Errorf("%w string is type %T", ErrInvalidInput, s)
}

func getUUID(u interface{}) (uuid.UUID, error) {
	switch u.(type) {
	case string:
		id, err := uuid.Parse(u.(string))
		if err != nil {
			return uuid.Nil, fmt.Errorf("%w UUID is %s", ErrInvalidInput, u.(string))
		}
		return id, nil
	case uuid.UUID:
		return u.(uuid.UUID), nil
	}
	return uuid.Nil, fmt.Errorf("%w UUID is type %T", ErrInvalidInput, u)
}

func getModel(mn interface{}, tx db.RWTx) (db.Model, error) {
	var out db.Model
	name, err := getString(mn)
	if err != nil {
		return out, err
	}
	return tx.GetModel(name)
}

func getMatcher(mm interface{}) (db.Matcher, error) {
	if val, ok := mm.(db.Matcher); ok {
		return val, nil
	}
	return nil, fmt.Errorf("%w Matcher is type %T", ErrInvalidInput, mm)
}

//Wrapper for the Record interface so we can control which methods to expose.
// This gets surfaced in Starlark as return values of database functions
type Record interface {
	ID() uuid.UUID
	Get(string) (interface{}, error)
	GetFK(string) (uuid.UUID, error)
}

type starlarkRecord struct {
	inner db.Record
}

func (r *starlarkRecord) ID() uuid.UUID {
	return r.inner.ID()
}

func (r *starlarkRecord) Get(fieldName string) (interface{}, error) {
	field, err := r.inner.Get(fieldName)
	if err != nil {
		return nil, err
	}
	return field, nil
}

func (r *starlarkRecord) GetFK(fieldName string) (uuid.UUID, error) {
	rel, err := r.inner.GetFK(fieldName)
	if err != nil {
		return uuid.Nil, err
	}
	return rel, nil
}

//Actual DB API
func DBLib(tx db.RWTx) map[string]interface{} {
	env := map[string]interface{}{
		"FindOne": func(mn, mm interface{}) (Record, error) {
			m, err := getModel(mn, tx)
			if err != nil {
				return nil, err
			}
			ma, err := getMatcher(mm)
			if err != nil {
				return nil, err
			}
			r, err := tx.FindOne(m.ID, ma)
			if err != nil {
				return nil, err
			}
			return &starlarkRecord{inner: r}, nil
		},
		"FindMany": func(mn, mm interface{}) ([]Record, error) {
			m, err := getModel(mn, tx)
			if err != nil {
				return nil, err
			}
			ma, err := getMatcher(mm)
			if err != nil {
				return nil, err
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
		},
		"Eq": func(k, v interface{}) (db.Matcher, error) {
			key, err := getString(k)
			if err != nil {
				return nil, err
			}
			return db.Eq(key, v), nil
		},
		"EqFK": func(k, v interface{}) (db.Matcher, error) {
			key, err := getString(k)
			if err != nil {
				return nil, err
			}
			id, err := getUUID(v)
			if err != nil {
				return nil, err
			}
			return db.EqFK(key, id), nil
		},
		"And": func(matchers ...interface{}) (db.Matcher, error) {
			var out []db.Matcher
			for i := 0; i < len(matchers); i++ {
				m, err := getMatcher(matchers[i])
				if err != nil {
					return nil, err
				}
				out = append(out, m)
			}
			return db.And(out...), nil
		},
		"Insert": func(mn interface{}, fields interface{}) (Record, error) {
			m, err := getModel(mn, tx)
			if err != nil {
				return nil, err
			}
			r := tx.MakeRecord(m.ID)
			err = r.Set("id", uuid.New())
			if err != nil {
				return nil, err
			}
			if fieldMap, ok := fields.(map[interface{}]interface{}); ok {
				for key, val := range fieldMap {
					ks, err := getString(key)
					if err != nil {
						return nil, err
					}
					err = r.Set(ks, recursiveFromValue(val.(starlark.Value)))
					if err != nil {
						return nil, err
					}
				}
			}
			tx.Insert(r)
			return &starlarkRecord{inner: r}, nil
		},
		"Update": func(r interface{}, fields interface{}) (Record, error) {
			if _, ok := r.(*starlarkRecord); !ok {
				return nil, fmt.Errorf("%w Type %T", ErrInvalidInput, r)
			}
			oldRec := r.(*starlarkRecord).inner
			newRec := oldRec.DeepCopy()
			if fieldMap, ok := fields.(map[interface{}]interface{}); ok {
				for key, val := range fieldMap {
					ks, err := getString(key)
					if err != nil {
						return nil, err
					}
					newRec.Set(ks, recursiveFromValue(val.(starlark.Value)))
				}
			}
			err := tx.Update(oldRec, newRec)
			if err != nil {
				return nil, err
			}
			return &starlarkRecord{inner: newRec}, err

		},
		"Connect": func(s interface{}, r1 interface{}, r2 interface{}) (bool, error) {
			bname, err := getString(s)
			if err != nil {
				return false, err
			}
			if _, ok := r1.(*starlarkRecord); !ok {
				return false, fmt.Errorf("%w Type %T", ErrInvalidInput, r1)
			}
			rec1 := r1.(*starlarkRecord)
			binding, err := rec1.inner.Model().GetBinding(bname)
			if err != nil {
				return false, err
			}
			if _, ok := r2.(*starlarkRecord); !ok {
				return false, fmt.Errorf("%w Type %T", ErrInvalidInput, r2)
			}
			rec2 := r2.(*starlarkRecord)
			err = tx.Connect(rec1.inner, rec2.inner, binding.Relationship)
			if err != nil {
				return false, err
			}
			return true, nil
		},
	}
	return env
}
