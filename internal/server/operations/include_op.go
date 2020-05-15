package operations

import (
	"awans.org/aft/internal/db"
	"encoding/json"
)

type Inclusion struct {
	Relationship db.Relationship
	Query        Query
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

// TODO hack -- remove this and rewrite with Relationship containing the name
func getBackref(tx db.Tx, rel db.Relationship) db.Relationship {
	m, _ := tx.GetModel(rel.TargetModel)
	return m.Relationships[rel.TargetRel]
}

func resolve(tx db.Tx, ir *IncludeResult, i Inclusion) error {
	rec := ir.Record
	id := ir.Record.Id()
	rel := i.Relationship
	backRel := getBackref(tx, rel)

	switch rel.RelType {
	case db.HasOne:
		// FK on the other side
		targetFK := db.JsonKeyToRelFieldName(rel.TargetRel)
		hit, err := tx.FindOne(rel.TargetModel, targetFK, id)
		if err != nil {
			return err
		}
		ir.SingleIncludes[backRel.TargetRel] = hit
	case db.BelongsTo:
		// FK on this side
		thisFK := rec.GetFK(backRel.TargetRel)
		hit, err := tx.FindOne(rel.TargetModel, "id", thisFK)
		if err != nil {
			return err
		}
		ir.SingleIncludes[backRel.TargetRel] = hit
	case db.HasMany:
		// FK on the other side
		hits := tx.FindMany(rel.TargetModel, db.EqFK(rel.TargetRel, id))
		ir.MultiIncludes[backRel.TargetRel] = hits
	case db.HasManyAndBelongsToMany:
		panic("Not implemented")
	}
	return nil
}
