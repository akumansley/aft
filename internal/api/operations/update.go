package operations

import (
	"fmt"

	"awans.org/aft/internal/db"
)

func (op UpdateOperation) Apply(tx db.RWTx) (*db.QueryResult, error) {
	root := tx.Ref(op.ModelID)
	clauses := HandleWhere(tx, root, op.FindArgs.Where)
	q := tx.Query(root, clauses...)
	outs := q.All()
	if len(outs) > 1 {
		return nil, fmt.Errorf("Found more than one record")
	}
	if len(outs) == 0 {
		return nil, fmt.Errorf("No record found")
	}
	for _, no := range op.Nested {
		err := no.ApplyNested(tx, root, outs)
		if err != nil {
			fmt.Printf("Hello\n")
			return nil, err
		}
	}
	oldRec := outs[0]
	newRec, err := updateRecordFromData(tx, oldRec.Record, op.Data)
	if err != nil {
		fmt.Printf("h1\n")
		return nil, err
	}
	err = tx.Update(oldRec.Record, newRec)
	if err != nil {
		fmt.Printf("h0\n")
		return nil, err
	}

	//rerun the query ensuring the right record is at the root
	root = tx.Ref(op.ModelID)
	clauses = []db.QueryClause{db.Filter(root, db.EqID(newRec.ID()))}
	clauses = append(clauses, handleIncludes(tx, root, op.FindArgs.Include)...)
	clauses = append(clauses, handleSelects(tx, root, op.FindArgs.Select)...)
	q = tx.Query(root, clauses...)
	outs = q.All()
	if len(outs) != 1 {
		return nil, fmt.Errorf("Resolve single include returned non-1 results")
	}
	fmt.Printf("h3\n")
	return outs[0], err
}

func (op NestedUpdateOperation) ApplyNested(tx db.RWTx, parent db.ModelRef, parents []*db.QueryResult) (err error) {
	outs, child := handleRelationshipWhere(tx, parent, parents, op.Relationship, op.Where)

	if len(outs) > 1 {
		return fmt.Errorf("Found more than one record")
	} else if len(outs) == 1 {
		for _, no := range op.Nested {
			err := no.ApplyNested(tx, child, outs)
			if err != nil {
				return err
			}
		}
		oldRec := outs[0].Record
		newRec, err := updateRecordFromData(tx, oldRec, op.Data)
		err = tx.Update(oldRec, newRec)
		outs[0].Record = newRec
		if err != nil {
			return err
		}
	}
	return
}

func updateRecordFromData(tx db.RWTx, oldRec db.Record, data map[string]interface{}) (db.Record, error) {
	newRec := oldRec.DeepCopy()
	m, err := tx.Schema().GetInterfaceByID(oldRec.InterfaceID())
	if err != nil {
		fmt.Printf("this is it\n")
		return nil, err
	}

	for k, v := range data {
		a, err := m.AttributeByName(tx, k)
		if err != nil {
			fmt.Printf("this 1 it\n")
			return nil, err
		}
		err = a.Set(tx, newRec, v)
		if err != nil {
			fmt.Printf("this 2 it %v\n", err)
			return nil, err
		}
		delete(data, k)
	}
	if len(data) != 0 {
		return nil, fmt.Errorf("Unused data in update")
	}
	return newRec, nil
}
