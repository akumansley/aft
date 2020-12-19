package operations

import (
	"fmt"

	"awans.org/aft/internal/db"
)

func (op CreateOperation) Apply(tx db.RWTx) (*db.QueryResult, error) {
	rec, err := buildRecordFromData(tx, op.ModelID, op.Data)
	if err != nil {
		return nil, err
	}
	tx.Insert(rec)

	root := tx.Ref(rec.Interface().ID())
	parents := []*db.QueryResult{&db.QueryResult{Record: rec}}
	for _, no := range op.Nested {
		err := no.ApplyNested(tx, root, parents)
		if err != nil {
			return nil, err
		}
	}

	ids := []db.ID{rec.ID()}
	clauses := []db.QueryClause{db.Filter(root, db.IDIn(ids))}
	clauses = append(clauses, handleIncludes(tx, root, op.FindArgs.Include)...)
	clauses = append(clauses, handleSelects(tx, root, op.FindArgs.Select)...)
	q := tx.Query(root, clauses...)
	qrs := q.All()
	if len(qrs) != 1 {
		return nil, fmt.Errorf("Resolve single include returned non-1 results")
	}
	return qrs[0], nil
}

func (op NestedCreateOperation) ApplyNested(tx db.RWTx, parent db.ModelRef, parents []*db.QueryResult) (err error) {
	rec, err := buildRecordFromData(tx, op.Model.ID(), op.Data)
	if err != nil {
		return err
	}
	tx.Insert(rec)
	for _, parent := range parents {
		tx.Connect(parent.Record.ID(), rec.ID(), op.Relationship.ID())
	}
	for _, no := range op.Nested {
		err = no.ApplyNested(tx, parent, []*db.QueryResult{&db.QueryResult{Record: rec}})
		if err != nil {
			return
		}
	}

	return nil
}

func (op NestedConnectOperation) ApplyNested(tx db.RWTx, parent db.ModelRef, parents []*db.QueryResult) (err error) {
	child := tx.Ref(op.Relationship.Target().ID())
	clauses := HandleWhere(tx, child, op.Where)
	q := tx.Query(child, clauses...)
	outs := q.All()

	if len(outs) > 1 {
		return fmt.Errorf("Found more than one record")
	} else if len(outs) == 1 {
		rec := outs[0].Record
		for _, parent := range parents {
			err = op.Relationship.Connect(parent.Record, rec)
			if err != nil {
				return err
			}
		}
	} else if len(outs) == 0 {
		return fmt.Errorf("No record found to connect")
	}
	return
}

func (op NestedDisconnectOperation) ApplyNested(tx db.RWTx, parent db.ModelRef, parents []*db.QueryResult) (err error) {
	outs, _ := handleRelationshipWhere(tx, parent, parents, op.Relationship, op.Where)

	if len(outs) > 1 {
		return fmt.Errorf("Found more than one record")
	} else if len(outs) == 1 {
		rec := outs[0].Record
		for _, parent := range parents {
			err = op.Relationship.Disconnect(parent.Record, rec)
			if err != nil {
				return err
			}
		}
	}
	return
}

func (op NestedSetOperation) ApplyNested(tx db.RWTx, parent db.ModelRef, parents []*db.QueryResult) (err error) {
	child := tx.Ref(op.Relationship.Target().ID())
	clauses := HandleWhere(tx, child, op.Where)
	q := tx.Query(child, clauses...)
	outs := q.All()

	if len(outs) > 1 {
		return fmt.Errorf("Found more than one record")
	} else if len(outs) == 1 {
		for _, parent := range parents {
			prec := parent.Record
			rec := outs[0].Record

			//disconnect the old stuff
			if op.Relationship.Multi() {
				olds, _ := op.Relationship.LoadMany(prec)
				for _, old := range olds {
					err = tx.Disconnect(prec.ID(), old.ID(), op.Relationship.ID())
					if err != nil {
						return err
					}
				}
			} else {
				old, err := op.Relationship.LoadOne(prec)
				if err == nil {
					err = tx.Disconnect(prec.ID(), old.ID(), op.Relationship.ID())
					if err != nil {
						return err
					}
				}
			}
			//connect the new stuff
			err = tx.Connect(prec.ID(), rec.ID(), op.Relationship.ID())
			if err != nil {
				return err
			}
		}
	} else if len(outs) == 0 {
		return fmt.Errorf("Tried to set to non-existant record")
	}
	return
}

func buildRecordFromData(tx db.RWTx, modelID db.ID, data map[string]interface{}) (db.Record, error) {
	m, err := tx.Schema().GetInterfaceByID(modelID)
	if err != nil {
		return nil, err
	}
	rec := db.NewRecord(m)
	for k, v := range data {
		a, err := m.AttributeByName(k)
		if err != nil {
			return nil, err
		}
		err = a.Set(rec, v)
		if err != nil {
			return nil, err
		}
	}
	return rec, nil
}
