package operations

import (
	"awans.org/aft/internal/db"
	"fmt"
)

func (op UpdateOperation) Apply(tx db.RWTx) (*db.QueryResult, error) {
	fo := FindOneOperation{ModelID: op.ModelID, Where: op.Where}
	oldRec, err := fo.Apply(tx)
	if err != nil {
		return nil, err
	}
	if oldRec == nil {
		return nil, fmt.Errorf("Can't find record to update")
	}
	newRec, err := updateRecordFromData(oldRec.Record, op.Data)
	if err != nil {
		return nil, err
	}
	err = tx.Update(oldRec.Record, newRec)
	if err != nil {
		return nil, err
	}
	outs[0].Record = newRec
	for _, no := range op.Nested {
		err := no.ApplyNested(tx)
		if err != nil {
			return nil, err
		}
	}
	err = tx.Commit()
	if err != nil {
		return nil, err
	}
	inc, err := op.Include.One(tx, oldRec.Record.Interface().ID(), newRec)
	return inc, err
}

func (op NestedUpdateOperation) ApplyNested(tx db.RWTx) (err error) {
	parent := tx.Ref(op.Relationship.Source().ID())
	child := tx.Ref(op.Relationship.Target().ID())

	root := tx.Ref(op.Relationship.Target().ID())
	clauses := handleWhere(tx, root, op.Where)
	on := parent.Rel(op.Relationship)
	clauses = append(clauses, db.Join(child, on))
	q := tx.Query(root, clauses...)
	out := q.All()

	if len(out) > 1 {
		return fmt.Errorf("Found more than one record")
	} else if len(out) == 1 {
		oldRec := out[0].Record
		newRec, err := updateRecordFromData(oldRec, op.Data)
		err = tx.Update(oldRec, newRec)
		if err != nil {
			return err
		}
		for _, no := range op.Nested {
			err := no.ApplyNested(tx)
			if err != nil {
				return err
			}
		}
	}
	return
}

func updateRecordFromData(oldRec db.Record, data map[string]interface{}) (db.Record, error) {
	newRec := oldRec.DeepCopy()
	attrs, err := oldRec.Interface().Attributes()
	if err != nil {
		return nil, err
	}
	for _, attr := range attrs {
		key := attr.Name()
		if value, ok := data[key]; ok {
			err = newRec.Set(key, value)
			if err != nil {
				return nil, err
			}
			delete(data, key)
		}
	}
	if len(data) != 0 {
		return nil, fmt.Errorf("Unused data in update")
	}
	return newRec, nil
}
