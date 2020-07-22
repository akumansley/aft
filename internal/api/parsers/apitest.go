package parsers

import (
	"awans.org/aft/internal/api/operations"
	"awans.org/aft/internal/db"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func CmpOpts() []cmp.Option {
	tFC := cmpopts.SortSlices(func(a, b operations.FieldCriterion) bool {
		return a.Key < b.Key
	})
	tRC := cmpopts.SortSlices(func(a, b operations.RelationshipCriterion) bool {
		return a.Relationship.Name() < b.Relationship.Name()
	})
	tA := cmpopts.SortSlices(func(a, b db.Attribute) bool {
		return a.Name() < b.Name()
	})
	ignoreFunc := cmpopts.IgnoreFields(db.NativeFunctionL{}, "Function")

	cmpModel := cmp.Comparer(func(a, b db.Model) bool {
		if a == nil || b == nil {
			return false
		}
		return a.ID() == b.ID()
	})
	cmpRel := cmp.Comparer(func(a, b db.Relationship) bool {
		if a == nil || b == nil {
			return false
		}
		return a.ID() == b.ID()
	})
	cmpAttr := cmp.Comparer(func(a, b db.Attribute) bool {
		if a == nil || b == nil {
			return false
		}
		return a.ID() == b.ID()
	})
	cmpDt := cmp.Comparer(func(a, b db.Datatype) bool {
		if a == nil || b == nil {
			return false
		}
		return a.ID() == b.ID()
	})
	return []cmp.Option{
		tFC,
		tRC,
		tA,
		ignoreFunc,
		cmpModel,
		cmpRel,
		cmpAttr,
		cmpDt,
	}
}

var IgnoreRecIDs = cmp.Comparer(func(a, b db.Record) bool {
	if a == nil || b == nil {
		return false
	}
	am := a.Map()
	match := true
	for k, av := range am {
		if k == "id" {
			continue
		}
		bv, _ := b.Get(k)
		if bv != av {
			match = false
			break
		}
	}
	return match
})
