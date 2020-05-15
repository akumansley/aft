package db

import (
	"awans.org/aft/internal/hold"
	"awans.org/aft/internal/model"
)

func (i Include) Resolve(tx Tx, rec model.Record) model.IncludeResult {
	ir := model.IncludeResult{Record: rec, SingleIncludes: make(map[string]model.Record), MultiIncludes: make(map[string][]model.Record)}

	for _, inc := range i.Includes {
		resolve(tx, &ir, inc)
	}
	return ir
}

func resolve(tx Tx, ir *model.IncludeResult, i Inclusion) error {
	rec := ir.Record
	id := ir.Record.Id()
	rel := i.Relationship
	backRel := getBackref(tx, rel)

	switch rel.RelType {
	case model.HasOne:
		// FK on the other side
		targetFK := model.JsonKeyToRelFieldName(rel.TargetRel)
		hit, err := tx.FindOne(rel.TargetModel, targetFK, id)
		if err != nil {
			return err
		}
		ir.SingleIncludes[backRel.TargetRel] = hit
	case model.BelongsTo:
		// FK on this side
		thisFK := rec.GetFK(backRel.TargetRel)
		hit, err := tx.FindOne(rel.TargetModel, "id", thisFK)
		if err != nil {
			return err
		}
		ir.SingleIncludes[backRel.TargetRel] = hit
	case model.HasMany:
		// FK on the other side
		hits := tx.FindMany(rel.TargetModel, hold.EqFK(rel.TargetRel, id))
		ir.MultiIncludes[backRel.TargetRel] = hits
	case model.HasManyAndBelongsToMany:
		panic("Not implemented")
	}
	return nil
}
