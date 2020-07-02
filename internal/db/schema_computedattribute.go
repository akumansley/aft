package db

import (
	"errors"
)

var ErrNotSupported = errors.New("set operation not supported on computed attributes")

type computedAttr struct {
	rec Record
	tx  Tx
}

func (a *computedAttr) ID() ID {
	return a.rec.ID()
}

func (a *computedAttr) Name() string {
	return caName.AsAttribute().MustGet(a.rec).(string)
}

func (a *computedAttr) Storage() EnumValue {
	return NotStored.AsEnumValue()
}

type GetterArgs struct {
	rec Record
	tx  Tx
}

func (a *computedAttr) Get(rec Record) (interface{}, error) {
	getterRel, _ := a.tx.Schema().GetRelationshipByID(ComputedAttributeGetter.ID)
	getter, err := a.tx.GetRelatedOne(a.ID(), getterRel)
	if err != nil {
		return nil, err
	}
	f, err := a.tx.loadFunction(getter)
	if err != nil {
		return nil, err
	}

	// do we need different args? maybe
	res, err := f.Call(GetterArgs{rec, a.tx})
	return res, err
}

func (a *computedAttr) MustGet(rec Record) interface{} {
	v, err := a.Get(rec)
	if err != nil {
		panic(err)
	}
	return v
}

func (a *computedAttr) Set(v interface{}, rec Record) error {
	return ErrNotSupported
}
