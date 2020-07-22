package operations

import (
	"awans.org/aft/internal/db"
	"fmt"
)

func (op FindOneOperation) Apply(tx db.Tx) (*db.QueryResult, error) {
	fm := FindManyOperation{ModelID: op.ModelID, Where: op.Where, Include: op.Include}
	out, err := fm.Apply(tx)
	if err != nil {
		return nil, err
	}
	if len(out) > 1 {
		return nil, fmt.Errorf("Found more than one record")
	}
	if len(out) == 0 {
		return nil, nil
	}
	return out[0], nil
}
