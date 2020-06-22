package api

import (
	"awans.org/aft/internal/db"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func CmpOpts() []cmp.Option {
	tFC := cmpopts.SortSlices(func(a, b FieldCriterion) bool {
		return a.Key < b.Key
	})
	tRC := cmpopts.SortSlices(func(a, b RelationshipCriterion) bool {
		return a.Relationship.Name < b.Relationship.Name
	})
	tA := cmpopts.SortSlices(func(a, b db.Attribute) bool {
		return a.Name < b.Name
	})
	ignoreFunc := cmpopts.IgnoreFields(db.Code{}, "Function")
	return []cmp.Option{
		tFC,
		tRC,
		tA,
		ignoreFunc,
	}
}
