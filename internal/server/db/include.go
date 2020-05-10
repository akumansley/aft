package db

import (
	"awans.org/aft/internal/model"
)

func (i Include) Resolve(tx Tx, rec model.Record) model.IncludeResult {
	ir := model.IncludeResult{Record: rec, SingleIncludes: make(map[string]model.Record), MultiIncludes: make(map[string][]model.Record)}
	for _, inc := range i.Includes {
		tx.Resolve(&ir, inc)
	}
	return ir
}
