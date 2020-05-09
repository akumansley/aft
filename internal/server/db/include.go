package db

import (
	"awans.org/aft/internal/model"
)

// TODO maybe RecWithInclude
func (i Include) Resolve(tx Tx, rec model.Record) model.Record {
	for _, inc := range i.Includes {
		tx.Resolve(rec, inc)
	}
	return rec
}
