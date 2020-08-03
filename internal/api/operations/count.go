package operations

import (
	"awans.org/aft/internal/db"
)

func (op CountOperation) Apply(tx db.Tx) (int, error) {
	root := tx.Ref(op.ModelID)
	q := tx.Query(root, HandleWhere(tx, root, op.Where)...)
	return len(q.All()), nil
}
