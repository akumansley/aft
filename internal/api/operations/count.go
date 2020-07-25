package operations

import (
	"awans.org/aft/internal/db"
)

func (op CountOperation) Apply(tx db.Tx) (int, error) {
	fm := FindManyOperation{ModelID: op.ModelID, Where: op.Where}
	q := fm.handleFindMany(tx)
	qrs := q.All()
	results := []db.Record{}
	for _, qr := range qrs {
		results = append(results, qr.Record)
	}

	return len(results), nil
}
