package api

import (
	"awans.org/aft/internal/db"
)

type Inclusion struct {
	Binding db.Binding
	Where   Where
}

type Include struct {
	Includes []Inclusion
}

func (i Include) Resolve(tx db.Tx, rec db.Record) db.QueryResult {
	ir := db.QueryResult{Record: rec, SingleRelated: make(map[string]db.Record), MultiRelated: make(map[string][]db.Record)}

	for _, inc := range i.Includes {
		resolve(tx, &ir, inc)
	}
	return ir
}

func resolve(tx db.Tx, ir *db.QueryResult, i Inclusion) error {
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
		ir.SingleRelated[b.Name()] = hit
	case db.BelongsTo:
		// FK on this side
		thisFK := rec.GetFK(b.Name())
		hit, err := tx.FindOne(d.ModelID(), db.Eq("id", thisFK))
		if err != nil {
			return err
		}
		ir.SingleRelated[b.Name()] = hit
	case db.HasMany:
		// FK on the other side
		hits := tx.FindMany(d.ModelID(), db.EqFK(d.Name(), id))
		ir.MultiRelated[b.Name()] = hits
	case db.HasManyAndBelongsToMany:
		panic("Not implemented")
	}
	return nil
}
