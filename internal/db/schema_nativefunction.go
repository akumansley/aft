package db

import (
	"github.com/google/uuid"
)

type nativeFunction struct {
	rec Record
	tx  Tx
}

func (nf nativeFunction) ID() ID {
	return nf.rec.ID()
}

func (nf nativeFunction) Name() string {
	return nf.rec.MustGet("name").(string)
}

func (nf nativeFunction) Storage() StorageEnumValue {
	_ = nf.rec.MustGet("storedAs").(uuid.UUID)
	panic("Not Implemented")
}

func (nf nativeFunction) FromJSON() Function {
	panic("Not Implemented")
}
