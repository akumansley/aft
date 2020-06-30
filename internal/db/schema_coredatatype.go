package db

import (
	"github.com/google/uuid"
)

type coreDatatype struct {
	rec Record
	tx  Tx
}

func (cd coreDatatype) ID() ID {
	return cd.rec.ID()
}

func (cd coreDatatype) Name() string {
	return cd.rec.MustGet("name").(string)
}

func (cd coreDatatype) Storage() StorageEnumValue {
	_ = cd.rec.MustGet("storedAs").(uuid.UUID)
	panic("Not Implemented")
}

func (cd coreDatatype) FromJSON() Function {
	panic("Not Implemented")
}

// func (d coreDatatype) FromJSON(arg interface{}) (interface{}, error) {
// 	c := d.Validator
// 	out, err := c.Executor.Invoke(c, arg)
// 	if err != nil {
// 		panic(err)
// 	}
// 	return typeCheck(d, out)
// }

// func (d coreDatatype) GetID() ID {
// 	return d.ID
// }

// func (d coreDatatype) Storage() StorageEnumValue {
// 	return d.StoredAs
// }

// func (d coreDatatype) FillRecord(storeDatatype Record) error {
// 	ew := NewRecordWriter(storeDatatype)
// 	ew.Set("id", uuid.UUID(d.ID))
// 	ew.Set("name", d.Name)
// 	ew.Set("storedAs", uuid.UUID(d.StoredAs.ID))
// 	ew.Set("enum", false)
// 	ew.Set("native", true)
// 	return ew.err
// }
