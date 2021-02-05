package operations

import (
	"fmt"

	"awans.org/aft/internal/db"
)

func (op NestedConnectOperation) ApplyNested(tx db.RWTx, parent db.ModelRef, parents []*db.QueryResult) (err error) {
	child := tx.Ref(op.Relationship.Target(tx).ID())
	clauses := HandleWhere(tx, child, op.Where)
	q := tx.Query(child, clauses...)
	outs := q.All()

	if len(outs) > 1 {
		return fmt.Errorf("Found more than one record")
	} else if len(outs) == 1 {
		rec := outs[0].Record
		for _, parent := range parents {
			err = op.Relationship.Connect(tx, parent.Record, rec)
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
			err = op.Relationship.Disconnect(tx, parent.Record, rec)
			if err != nil {
				return err
			}
		}
	}
	return
}

func (op NestedSetOperation) ApplyNested(tx db.RWTx, parent db.ModelRef, parents []*db.QueryResult) (err error) {
	child := tx.Ref(op.Relationship.Target(tx).ID())
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
				olds, _ := op.Relationship.LoadMany(tx, prec)
				for _, old := range olds {
					err = tx.Disconnect(prec.ID(), old.ID(), op.Relationship.ID())
					if err != nil {
						return err
					}
				}
			} else {
				old, err := op.Relationship.LoadOne(tx, prec)
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
