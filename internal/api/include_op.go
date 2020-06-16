package api

import (
	"awans.org/aft/internal/db"
	"encoding/json"
)

type Inclusion struct {
	Binding db.Binding
	Where   Where
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
	rec := ir.Record
	id := ir.Record.ID()
	b := i.Binding
	d := b.Dual()

	switch b.RelType() {
	case db.HasOne:
		// FK on the other side
		hit, err := tx.FindOne(d.ModelID(), db.EqFK(d.Name(), id))
		if err != nil {
			return err
		}
		ir.SingleIncludes[b.Name()] = hit
	case db.BelongsTo:
		// FK on this side
		thisFK := rec.GetFK(b.Name())
		hit, err := tx.FindOne(d.ModelID(), db.Eq("id", thisFK))
		if err != nil {
			return err
		}
		ir.SingleIncludes[b.Name()] = hit
	case db.HasMany:
		// FK on the other side
		hits := tx.FindMany(d.ModelID(), db.EqFK(d.Name(), id))
		ir.MultiIncludes[b.Name()] = hits
	case db.HasManyAndBelongsToMany:
		panic("Not implemented")
	}
	return nil
}
