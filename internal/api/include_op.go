package api

import (
	"awans.org/aft/internal/db"
	"encoding/json"
)

type Inclusion struct {
	Relationship db.Relationship
	Where        Where
}

type Include struct {
	Includes []Inclusion
}

type IncludeResult struct {
	Record         db.Record
	SingleIncludes map[string]db.Record
	MultiIncludes  map[string][]db.Record
}

func (ir IncludeResult) MarshalJSON() ([]byte, error) {
	data := ir.Record.Map()
	for k, v := range ir.SingleIncludes {
		data[k] = v
	}
	for k, v := range ir.MultiIncludes {
		data[k] = v
	}
	return json.Marshal(data)
}

func (i Include) Resolve(tx db.Tx, rec db.Record) IncludeResult {
	ir := IncludeResult{Record: rec, SingleIncludes: make(map[string]db.Record), MultiIncludes: make(map[string][]db.Record)}

	for _, inc := range i.Includes {
		resolve(tx, &ir, inc)
	}
	return ir
}

func resolve(tx db.Tx, ir *IncludeResult, i Inclusion) error {
	// rec := ir.Record
	// id := ir.Record.ID()
	// r := i.Relationship
	// TODO: rewrite on query interface
	panic("Not Implemented")

}
